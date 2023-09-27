package client

import (
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
)

const (
	idxTimestamp = iota
	idxMissedBlockNumbers
	idxLastSyncedBlockNumber
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
		name: "missed_block_numbers",
		desc: "TODO",
	},
	{
		name: "last_synced_block_number",
		desc: "TODO",
	},
}

func initSyncsTable() *schema.Table {
	table := new(schema.Table)
	table.Name = "syncs"
	table.Description = "sync statistics"
	table.Columns = schema.ColumnList{
		schema.Column{
			Name:           syncColumnValues[idxTimestamp].name,
			Description:    syncColumnValues[idxTimestamp].desc,
			NotNull:        false,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.FixedWidthTypes.Timestamp_us,
		},
		schema.Column{
			Name:           syncColumnValues[idxMissedBlockNumbers].name,
			Description:    syncColumnValues[idxMissedBlockNumbers].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncColumnValues[idxLastSyncedBlockNumber].name,
			Description:    syncColumnValues[idxLastSyncedBlockNumber].desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
			//Type:           arrow.NumericFieldType, // TODO
		},
	}
	return table
}
