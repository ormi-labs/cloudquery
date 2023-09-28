package client

import (
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

const (
	idxTimestamp = iota
	idxStartBlock
	idxEnabledEntities
	idxMissedBlocks
	idxLastSyncedBlock
	idxEachSyncedBlock
	idxTimeElapsed
)

type syncColumn struct {
	name  string
	desc  string
	value any
}

var syncColumnValues = []syncColumn{
	{
		name: "timestamp",
		desc: "TODO",
	},
	{
		name: "start_block",
		desc: "TODO",
	},
	{
		name: "enabled_entities",
		desc: "TODO",
	},
	{
		name: "missed_blocks",
		desc: "TODO",
	},
	{
		name: "last_synced_block",
		desc: "TODO",
	},
	{
		name: "each_synced_block",
		desc: "TODO",
	},
	{
		name: "time_elapsed",
		desc: "TODO",
	},
}

func initSyncsTable() *schema.Table {
	table := new(schema.Table)
	table.Name = "syncs"
	table.Description = "sync run report"
	table.Columns = schema.ColumnList{
		schema.Column{
			Name:           syncColumnValues[idxTimestamp].name,
			Description:    syncColumnValues[idxTimestamp].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.FixedWidthTypes.Timestamp_us,
		},
		schema.Column{
			Name:           syncColumnValues[idxStartBlock].name,
			Description:    syncColumnValues[idxStartBlock].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.PrimitiveTypes.Uint64,
		},
		schema.Column{
			Name:           syncColumnValues[idxEnabledEntities].name,
			Description:    syncColumnValues[idxEnabledEntities].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncColumnValues[idxMissedBlocks].name,
			Description:    syncColumnValues[idxMissedBlocks].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncColumnValues[idxLastSyncedBlock].name,
			Description:    syncColumnValues[idxLastSyncedBlock].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           &arrow.Decimal128Type{},
		},
		schema.Column{
			Name:           syncColumnValues[idxEachSyncedBlock].name,
			Description:    syncColumnValues[idxEachSyncedBlock].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncColumnValues[idxTimeElapsed].name,
			Description:    syncColumnValues[idxTimeElapsed].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
	}
	return table
}
