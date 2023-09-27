package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudquery/cloudquery/plugins/source/postgresql/client/sql"
	"strings"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scalar"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type entity struct {
	tableName    string
	arrowSchema  *arrow.Schema
	recBuilder   *array.RecordBuilder
	transformers []transformer
	colNames     []string
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

	c.logger.Info().Any("entities", c.pluginSpec.Entities).Msg("enabled entities")

	names := make([]string, len(filteredTables))
	for i, table := range filteredTables {
		names[i] = table.Name
	}
	c.logger.Info().Strs("names", names).Msg("filtered tables")

	switch c.pluginSpec.SyncMode {
	case "table":
		for _, table := range filteredTables {
			if err := c.syncTable(ctx, tx, table, res); err != nil {
				return err
			}
		}
	case "block":
		// syncBlock is altered method to process blockchain data in orderly fashion
		if err := c.syncBlocks(ctx, tx, filteredTables, res); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit sync transaction: %w", err)
	}
	return nil
}

func (c *Client) query(ctx context.Context, entity *entity, sql string, args ...any) error {
	rows, err := c.Conn.Query(ctx, sql, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return err
		}

		for i, value := range values {
			val, err := entity.transformers[i](value)
			if err != nil {
				return err
			}

			s := scalar.NewScalar(entity.arrowSchema.Field(i).Type)
			if err := s.Set(val); err != nil {
				return err
			}

			scalar.AppendToBuilder(entity.recBuilder.Field(i), s)
		}
	}
	return nil
}

func (c *Client) syncBlocks(ctx context.Context, tx pgx.Tx, tables []*schema.Table, res chan<- message.SyncMessage) error {
	entities := make(map[EntityName]*entity)
	valid := false

LOOP:
	for _, e := range c.pluginSpec.Entities {
		if e.Enabled {
			for _, table := range tables {
				if table.Name == e.Table {
					// there has to be block entity
					if e.Name == EntBlock {
						valid = true
					}

					arrowSchema := table.ToArrowSchema()
					entities[e.Name] = &entity{
						tableName:    table.Name,
						arrowSchema:  arrowSchema,
						recBuilder:   array.NewRecordBuilder(memory.DefaultAllocator, arrowSchema),
						transformers: transformersForSchema(arrowSchema),
						colNames:     make([]string, len(table.Columns)),
					}
					for i, col := range table.Columns {
						entities[e.Name].colNames[i] = pgx.Identifier{col.Name}.Sanitize()
					}
					continue LOOP
				}
			}
			// no table
			return fmt.Errorf("no/wrong table specified for %q entity", e.Name)
		} else {
			entities[e.Name] = nil
		}
	}

	if !valid {
		return fmt.Errorf("no block entity specified in %q sync mode", "block")
	}

	var (
		syncsSchema  *arrow.Schema
		syncsBuilder *array.RecordBuilder
		//missed       any
		last_synced any
	)

	for _, table := range tables {
		if table.Name == "syncs" {
			syncsSchema = table.ToArrowSchema()
			syncsBuilder = array.NewRecordBuilder(memory.DefaultAllocator, syncsSchema)
		}
	}

	// block
	builder := new(sql.Select)
	blockSQL, err := builder.
		Select(entities[EntBlock].colNames...).
		From(pgx.Identifier{entities[EntBlock].tableName}.Sanitize()).
		Limit(int64(c.pluginSpec.Limit)).
		ToSQL()
	if err != nil {
		return fmt.Errorf("build query error: %v", err)
	}
	c.logger.Info().Str("SQL", blockSQL).Msg("block SQL")

	// transaction
	builder = new(sql.Select)
	txSQL, err := builder.
		Select(entities[EntTransaction].colNames...).
		From(pgx.Identifier{entities[EntTransaction].tableName}.Sanitize()).
		Where("block_hash", "=", "$1").
		ToSQL()
	if err != nil {
		return fmt.Errorf("build query error: %v", err)
	}
	c.logger.Info().Str("SQL", txSQL).Msg("transaction SQL")

	// TODO logs
	// TODO traces

	blocks, err := tx.Query(ctx, blockSQL)
	if err != nil {
		return err
	}
	defer blocks.Close()

	var (
		hashIdx   int
		numberIdx int
	)
	for i, fd := range blocks.FieldDescriptions() {
		switch fd.Name {
		case "hash":
			hashIdx = i
		case "number":
			numberIdx = i
		}
	}

	for blocks.Next() {
		values, err := blocks.Values()
		if err != nil {
			return err
		}

		for i, value := range values {
			val, err := entities[EntBlock].transformers[i](value)
			if err != nil {
				return err
			}

			s := scalar.NewScalar(entities[EntBlock].arrowSchema.Field(i).Type)
			if err := s.Set(val); err != nil {
				return err
			}

			scalar.AppendToBuilder(entities[EntBlock].recBuilder.Field(i), s)
		}

		c.logger.Info().Msg(values[hashIdx].(string))

		// transaction
		//c.query(ctx, entities[EntTransaction], txSQL)
		c.query(ctx, entities[EntTransaction], txSQL, values[hashIdx])
		res <- &message.SyncInsert{Record: entities[EntTransaction].recBuilder.NewRecord()}

		// NewRecord resets the builder for reuse
		res <- &message.SyncInsert{Record: entities[EntBlock].recBuilder.NewRecord()}

		last_synced = values[numberIdx]
	}

	syncColumnValues[idxTimestamp].value = time.Now()
	syncColumnValues[idxMissedBlockNumbers].value = "TODO"
	syncColumnValues[idxLastSyncedBlockNumber].value = fmt.Sprintf("%v", last_synced)
	if err := c.syncSync(res, syncsSchema, syncsBuilder, syncColumnValues); err != nil {
		return fmt.Errorf("sync sync error: %v", err)
	}

	return nil
}

func (c *Client) syncSync(res chan<- message.SyncMessage, schema *arrow.Schema, builder *array.RecordBuilder, vals []syncColumn) error {
	for i, val := range vals {
		s := scalar.NewScalar(schema.Field(i).Type)
		if err := s.Set(val.value); err != nil {
			return err
		}
		scalar.AppendToBuilder(builder.Field(i), s)
	}
	res <- &message.SyncInsert{Record: builder.NewRecord()}
	return nil
}

func (c *Client) syncTable(ctx context.Context, tx pgx.Tx, table *schema.Table, res chan<- message.SyncMessage) error {
	arrowSchema := table.ToArrowSchema()
	builder := array.NewRecordBuilder(memory.DefaultAllocator, arrowSchema)
	transformers := transformersForSchema(arrowSchema)

	colNames := make([]string, 0, len(table.Columns))
	for _, col := range table.Columns {
		colNames = append(colNames, pgx.Identifier{col.Name}.Sanitize())
	}
	query := "SELECT " + strings.Join(colNames, ",") + " FROM " + pgx.Identifier{table.Name}.Sanitize()
	rows, err := tx.Query(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	rowsInRecord := 0
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return err
		}

		for i, value := range values {
			val, err := transformers[i](value)
			if err != nil {
				return err
			}

			s := scalar.NewScalar(arrowSchema.Field(i).Type)
			if err := s.Set(val); err != nil {
				return err
			}

			scalar.AppendToBuilder(builder.Field(i), s)
		}

		rowsInRecord++
		if rowsInRecord >= c.pluginSpec.RowsPerRecord {
			res <- &message.SyncInsert{Record: builder.NewRecord()} // NewRecord resets the builder for reuse
			rowsInRecord = 0
		}
	}

	record := builder.NewRecord()
	if record.NumRows() > 0 { // only send if there are some unsent rows
		res <- &message.SyncInsert{Record: record}
	}

	return nil
}

func stringForTime(t pgtype.Time, unit arrow.TimeUnit) string {
	extra := ""
	hour := t.Microseconds / 1e6 / 60 / 60
	minute := t.Microseconds / 1e6 / 60 % 60
	second := t.Microseconds / 1e6 % 60
	micros := t.Microseconds % 1e6
	switch unit {
	case arrow.Millisecond:
		extra = fmt.Sprintf(".%03d", (micros)/1e3)
	case arrow.Microsecond:
		extra = fmt.Sprintf(".%06d", micros)
	case arrow.Nanosecond:
		// postgres doesn't support nanosecond precision
		extra = fmt.Sprintf(".%06d", micros)
	}

	return fmt.Sprintf("%02d:%02d:%02d"+extra, hour, minute, second)
}
