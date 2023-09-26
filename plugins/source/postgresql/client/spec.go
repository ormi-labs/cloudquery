package client

type EntityName string

const (
	EntBlock       EntityName = "block"
	EntTransaction            = "transaction"
	EntLog                    = "log"
	EntTrace                  = "trace"
)

type Entity struct {
	Name    EntityName `json:"name,omitempty"`
	Table   string     `json:"table,omitempty"`
	Enabled bool       `json:"enable,omitempty"`
}

type Spec struct {
	ConnectionString string   `json:"connection_string,omitempty"`
	PgxLogLevel      LogLevel `json:"pgx_log_level,omitempty"`
	CDCId            string   `json:"cdc_id,omitempty"`
	RowsPerRecord    int      `json:"rows_per_record,omitempty"`
	SyncMode         string   `json:"sync_mode,omitempty"`
	Limit            int      `json:"limit,omitempty"`
	Entities         []Entity `json:"entities,omitempty"`
}

func (s *Spec) SetDefaults() {
	if s.RowsPerRecord < 1 {
		s.RowsPerRecord = 1
	}
	if s.SyncMode == "" {
		s.SyncMode = "table"
	}
}
