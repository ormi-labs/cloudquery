package client

import (
	"encoding/json"
	"fmt"
	"github.com/apache/arrow/go/v14/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"os"
	"path/filepath"
)

const (
	idxSyncID = iota
	idxTimestamp
	idxStartBlock
	idxEnabledEntities
	idxMissedBlocks
	idxLastSyncedBlock
	idxEachSyncedBlock
	idxTimeElapsed
)

type syncColumns []struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Value any    `json:"value"`
}

var syncReport = syncColumns{
	idxSyncID: {
		Name: "sync_id",
		Desc: "TODO",
	},
	idxTimestamp: {
		Name: "timestamp",
		Desc: "TODO",
	},
	idxStartBlock: {
		Name: "start_block",
		Desc: "TODO",
	},
	idxEnabledEntities: {
		Name: "enabled_entities",
		Desc: "TODO",
	},
	idxMissedBlocks: {
		Name: "missed_blocks",
		Desc: "TODO",
	},
	idxLastSyncedBlock: {
		Name: "last_synced_block",
		Desc: "TODO",
	},
	idxEachSyncedBlock: {
		Name: "each_synced_block",
		Desc: "TODO",
	},
	idxTimeElapsed: {
		Name: "time_elapsed",
		Desc: "TODO",
	},
}

func initSyncsTable() *schema.Table {
	table := new(schema.Table)
	table.Name = "syncs"
	table.Description = "sync run report"
	table.Columns = schema.ColumnList{
		schema.Column{
			Name:           syncReport[idxSyncID].Name,
			Description:    syncReport[idxSyncID].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     true,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncReport[idxTimestamp].Name,
			Description:    syncReport[idxTimestamp].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.FixedWidthTypes.Timestamp_us,
		},
		schema.Column{
			Name:           syncReport[idxStartBlock].Name,
			Description:    syncReport[idxStartBlock].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.PrimitiveTypes.Uint64,
		},
		schema.Column{
			Name:           syncReport[idxEnabledEntities].Name,
			Description:    syncReport[idxEnabledEntities].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncReport[idxMissedBlocks].Name,
			Description:    syncReport[idxMissedBlocks].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncReport[idxLastSyncedBlock].Name,
			Description:    syncReport[idxLastSyncedBlock].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			//Type:           &arrow.FixedWidthTypes.Decimal128Type{},
			Type: arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncReport[idxEachSyncedBlock].Name,
			Description:    syncReport[idxEachSyncedBlock].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
		schema.Column{
			Name:           syncReport[idxTimeElapsed].Name,
			Description:    syncReport[idxTimeElapsed].Desc,
			NotNull:        true,
			IncrementalKey: false,
			PrimaryKey:     false,
			Type:           arrow.BinaryTypes.String,
		},
	}
	return table
}

func readFromFile(file string) (syncColumns, error) {
	_, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	// TODO

	return nil, nil
}

func (s syncColumns) writeToFile(dir string) error {
	f, err := os.OpenFile(filepath.Join(dir, s[idxSyncID].Value.(string)), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	bs, err := json.Marshal(s)
	if err != nil {
		fmt.Errorf("JSON marshal error: %v", err)
	}

	_, err = f.Write(bs)
	return err
}
