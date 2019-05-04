package storage

// Fields struct for storing fields data
type Fields struct {
	ID    int64  `reindex:"id,,pk"` // autoincrement basically
	Field string `reindex:"field"`  // event_id, src_ip and etc
}

// Events struct for storing events data
type Events struct {
	ID      int64                  `reindex:"id,,pk"`             // autoincrement basically
	LogName string                 `reindexer:"logname"`          // Security, System and etc
	EventID int32                  `reindexer:"eventid"`          // 1, 4638 and etc
	Data    map[string]interface{} `json:"data" reindexer:"data"` // '{"event_id":1,"event_log_name":"sysmon"}'
}

// Database represents an interface for a storage
type Database interface {
	Open(url string) error
	InsertFields(field *Fields) error
	InsertEvents(event *Events) error
	RestoreFields() ([]*Fields, error)
	RestoreEvents() ([]*Events, error)
	DeleteEvents(logname string, eventid int32) error
	DeleteFields(field string) error
	GetByEvent(logname string, eventid int32) ([]byte, error)
}
