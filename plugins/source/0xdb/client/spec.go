package client

type EntityName string

const (
	EntTransaction EntityName = "transaction"
	EntLog                    = "log"
	EntTrace                  = "trace"
)

type Block struct {
	Limit    int      `json:"limit,omitempty"`
	Start    uint64   `json:"start,omitempty"`
	TableIdx int      `json:"table_idx,omitempty"`
	Entities []Entity `json:"entities,omitempty"`
}

type Entity struct {
	Name     EntityName `json:"name,omitempty"`
	TableIdx int        `json:"table_idx,omitempty"`
	Enabled  bool       `json:"enable,omitempty"`
}

type Spec struct {
	ConnectionString string   `json:"connection_string,omitempty"`
	PgxLogLevel      LogLevel `json:"pgx_log_level,omitempty"`
	CDCId            string   `json:"cdc_id,omitempty"`
	RowsPerRecord    int      `json:"rows_per_record,omitempty"`
	ReportDir        string   `json:"report_dir,omitempty"`
	ReportFmt        string   `json:"report_fmt,omitempty"`
	Block            Block    `json:"block,omitempty"`
}

func (s *Spec) SetDefaults() {
	if s.RowsPerRecord < 1 {
		s.RowsPerRecord = 1
	}
	if s.ReportDir == "" {
		s.ReportDir = "."
	}
	if s.ReportFmt == "" {
		s.ReportFmt = "csv"
	}
}

func (en EntityName) String() string {
	return string(en)
}
