package client

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/apache/arrow/go/v14/arrow/decimal128"
	"github.com/cloudquery/cloudquery/plugins/source/postgresql/client/sql"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/jackc/pgx/v5"
)

type entity struct {
	tableName    string
	arrowSchema  *arrow.Schema
	recBuilder   *array.RecordBuilder
	transformers []transformer
	colNames     []string
	sqlStmt      string
}

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	if c.options.NoConnection {
		return fmt.Errorf("no connection")
	}
	var err error
	var snapshotName string

	connPool, err := c.Conn.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	// this must be closed only at the end of the initial sync process otherwise the snapshot
	// used to sync the initial data will be released
	defer connPool.Release()
	conn := connPool.Conn().PgConn()

	filteredTables, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	filteredTables = append(filteredTables, initSyncsTable())

	for _, table := range filteredTables {
		res <- &message.SyncMigrateTable{
			Table: table,
		}
	}

	if c.pluginSpec.CDCId != "" {
		snapshotName, err = c.startCDC(ctx, filteredTables, conn)
		if err != nil {
			return err
		}
	}

	if c.pluginSpec.CDCId != "" && snapshotName == "" {
		c.logger.Info().Msg("cdc is enabled but replication slot already exists, skipping initial sync")
	} else {
		if err := c.syncTables(ctx, snapshotName, filteredTables, res); err != nil {
			return err
		}
	}

	if c.pluginSpec.CDCId == "" {
		return nil
	}

	if err := c.listenCDC(ctx, res); err != nil {
		return fmt.Errorf("failed to listen to cdc: %w", err)
	}
	return nil
}

func (c *Client) syncTables(ctx context.Context, snapshotName string, filteredTables schema.Tables, res chan<- message.SyncMessage) error {
	tx, err := c.Conn.BeginTx(ctx, pgx.TxOptions{
		// this transaction is needed for us to take a snapshot and we need to close it only at the end of the initial sync
		// https://www.postgresql.org/docs/current/transaction-iso.html
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			if !errors.Is(err, pgx.ErrTxClosed) {
				c.logger.Error().Err(err).Msg("failed to rollback sync transaction")
			}
		}
	}()

	if snapshotName != "" {
		// if we use cdc we need to set the snapshot that was exported when we started the logical
		// replication stream
		if _, err := tx.Exec(ctx, "SET TRANSACTION SNAPSHOT '"+snapshotName+"'"); err != nil {
			return fmt.Errorf("failed to 'SET TRANSACTION SNAPSHOT %s': %w", snapshotName, err)
		}
	}

	c.logger.Info().Any("entities", c.pluginSpec.Block.Entities).Msg("enabled entities")

	names := make([]string, len(filteredTables))
	for i, table := range filteredTables {
		names[i] = table.Name
	}
	c.logger.Info().Strs("names", names).Msg("filtered tables")

	if err := c.syncBlocks(ctx, tx, filteredTables, res); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit sync transaction: %w", err)
	}
	return nil
}

// calcMissed calculates and formats missed blocks if
// one or many are missed out from a sequence.
// For ex.:
// input: 0, 1 | output: ""
// input: 0, 2 | output: "1"
// input: 0, 3 | output: "1-2"
// input: 0, 4 | output: "1-3"
func calcMissed(prev *uint64, next uint64) string {
	defer func() {
		*prev = next
	}()
	switch delta := next - *prev; {
	case delta == 1:
		return "" // Ok
	case delta == 2:
		return strconv.FormatUint(*prev+1, 10)
	case delta > 2:
		return strconv.FormatUint(*prev+1, 10) + "-" +
			strconv.FormatUint(next-1, 10)
	}
	return ""
}

// syncBlock syncs blockchain data block by block.
func (c *Client) syncBlocks(ctx context.Context, tx pgx.Tx, tables []*schema.Table, res chan<- message.SyncMessage) error {
	var (
		syncsSchema     *arrow.Schema
		syncsBuilder    *array.RecordBuilder
		lastSynced      uint64
		enabledEntities []string
		missedBlocks    string
		prevBlock       = c.pluginSpec.Block.Start
	)

	for i, table := range tables {
		c.logger.Info().Int("index", i).Str("name", tables[0].Name).Msg("table")
		if table.Name == "syncs" {
			syncsSchema = table.ToArrowSchema()
			syncsBuilder = array.NewRecordBuilder(memory.DefaultAllocator, syncsSchema)
		}
	}

	tableInit := func(table *schema.Table) (name string, schema *arrow.Schema,
		builder *array.RecordBuilder, transformers []transformer, colNames []string) {
		name = table.Name
		schema = table.ToArrowSchema()
		builder = array.NewRecordBuilder(memory.DefaultAllocator, schema)
		transformers = transformersForSchema(schema)
		colNames = make([]string, len(table.Columns))
		for i, col := range table.Columns {
			colNames[i] = pgx.Identifier{col.Name}.Sanitize()
		}
		return
	}

	blockName, blockSchema, blockBuilder, blockTransformers, blockColNames :=
		tableInit(tables[c.pluginSpec.Block.TableIdx])

	blockSQL, err := new(sql.Select).
		Select(blockColNames...).
		From(pgx.Identifier{blockName}.Sanitize()).
		Where("number", ">=", strconv.FormatUint(c.pluginSpec.Block.Start, 10)).
		Order("number").
		Limit(int64(c.pluginSpec.Block.Limit)).
		ToSQL()
	if err != nil {
		return fmt.Errorf("build query error: %v", err)
	}
	c.logger.Info().Str("SQL", blockSQL).Msg("block SQL")

	entities := make(map[EntityName]*entity)
	for _, e := range c.pluginSpec.Block.Entities {
		if e.Enabled {
			enabledEntities = append(enabledEntities, e.Name.String())
			table := tables[e.TableIdx]

			entities[e.Name] = &entity{}
			entities[e.Name].tableName,
				entities[e.Name].arrowSchema,
				entities[e.Name].recBuilder,
				entities[e.Name].transformers,
				entities[e.Name].colNames =
				tableInit(table)

			switch e.Name {
			case EntTransaction:
				builder := new(sql.Select)
				entities[e.Name].sqlStmt, err = builder.
					Select(entities[EntTransaction].colNames...).
					From(pgx.Identifier{entities[EntTransaction].tableName}.Sanitize()).
					Where("block_number", "=", "$1").
					ToSQL()
				if err != nil {
					return fmt.Errorf("build query error: %v", err)
				}
				c.logger.Info().Str("SQL", entities[e.Name].sqlStmt).Msg("transaction SQL")
			case EntLog:
				builder := new(sql.Select)
				entities[e.Name].sqlStmt, err = builder.
					Select(entities[EntLog].colNames...).
					From(pgx.Identifier{entities[EntLog].tableName}.Sanitize()).
					Where("block_number", "=", "$1").
					ToSQL()
				if err != nil {
					return fmt.Errorf("build query error: %v", err)
				}
				c.logger.Info().Str("SQL", entities[e.Name].sqlStmt).Msg("log SQL")
			case EntTrace:
				builder := new(sql.Select)
				entities[e.Name].sqlStmt, err = builder.
					Select(entities[EntTrace].colNames...).
					From(pgx.Identifier{entities[EntTrace].tableName}.Sanitize()).
					Where("block_number", "=", "$1").
					ToSQL()
				if err != nil {
					return fmt.Errorf("build query error: %v", err)
				}
				c.logger.Info().Str("SQL", entities[e.Name].sqlStmt).Msg("trace SQL")
			}
		}
	}

	syncReport[idxSyncID].Value = "sync_" + time.Now().Format("Jan-2-15:04:05")
	syncReport[idxTimestamp].Value = time.Now()
	syncReport[idxStartBlock].Value = c.pluginSpec.Block.Start
	syncReport[idxEnabledEntities].Value = strings.Join(enabledEntities, ",")

	defer func() {
		if err := c.syncSync(res, syncsSchema, syncsBuilder, syncReport); err != nil {
			c.logger.Error().Err(err).Msg("sync sync error")
		}
	}()

	qp := queryPool{
		entities: entities,
		conn:     c.Conn,
		in:       make(chan query, 1),
		out:      make(chan error, 1),
		res:      res,
		quit:     make(chan struct{}),
	}

	qp.start(ctx)
	defer qp.stop()

	blocks, err := tx.Query(ctx, blockSQL)
	if err != nil {
		return err
	}
	defer blocks.Close()

	var blockNumColIdx int
	for i, fd := range blocks.FieldDescriptions() {
		if fd.Name == "number" {
			blockNumColIdx = i
		}
	}

	for blocks.Next() {
		values, err := blocks.Values()
		if err != nil {
			return err
		}

		for i, value := range values {
			val, err := blockTransformers[i](value)
			if err != nil {
				return err
			}

			s := scalar.NewScalar(blockSchema.Field(i).Type)
			if err := s.Set(val); err != nil {
				return err
			}

			scalar.AppendToBuilder(blockBuilder.Field(i), s)

			if i == blockNumColIdx {
				switch v := val.(type) {
				case decimal128.Num:
					lastSynced = v.BigInt().Uint64()
				default:
					return fmt.Errorf("unimplemented type for block number: %T", v)
				}
			}
		}

		if v := calcMissed(&prevBlock, lastSynced); v != "" {
			if missedBlocks != "" {
				missedBlocks += ","
			}
			missedBlocks += v
		}

		if entities[EntTransaction] != nil {
			qp.add(query{
				entity: EntTransaction,
				arg1:   lastSynced,
			})
		}

		if entities[EntLog] != nil {
			qp.add(query{
				entity: EntLog,
				arg1:   lastSynced,
			})
		}

		if entities[EntTrace] != nil {
			qp.add(query{
				entity: EntTrace,
				arg1:   lastSynced,
			})
		}

		for range entities {
			if err := <-qp.out; err != nil {
				return fmt.Errorf("query error: %v", err)
			}
		}

		// NewRecord resets the builder for reuse
		res <- &message.SyncInsert{Record: blockBuilder.NewRecord()}

		start := syncReport[idxTimestamp].Value.(time.Time)
		syncReport[idxTimeElapsed].Value = FmtElapsedTime(start, time.Now(), " ", false)
		syncReport[idxLastSyncedBlock].Value = strconv.FormatUint(lastSynced, 10)
		syncReport[idxMissedBlocks].Value = missedBlocks
		syncReport.writeToFile(c.pluginSpec.ReportDir, c.pluginSpec.ReportFmt)
	}

	return nil
}

func (c *Client) syncSync(res chan<- message.SyncMessage, schema *arrow.Schema, builder *array.RecordBuilder, vals syncColumns) error {
	for i, val := range vals {
		s := scalar.NewScalar(schema.Field(i).Type)
		if err := s.Set(val.Value); err != nil {
			return err
		}
		scalar.AppendToBuilder(builder.Field(i), s)
	}
	res <- &message.SyncInsert{Record: builder.NewRecord()}
	return nil
}
