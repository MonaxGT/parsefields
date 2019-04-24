package storage

// Database represents an interface for a storage

type Fields struct {
	ID    int64  `reindex:"id,,pk"`
	Field string `reindex:"field"`
}

type Events struct {
	ID      int64                  `reindex:"id,,pk"`
	LogName string                 `reindexer:"logname"`
	EventID int32                  `reindexer:"eventid"`
	Data    map[string]interface{} `json:"data" reindexer:"data"`
}

type Database interface {
	Open(url string) error
	InsertFields(field *Fields) error
	InsertEvents(event *Events) error
	RestoreFields() ([]*Fields, error)
	RestoreEvents() ([]*Events, error)
	DeleteEvents(logname string, eventid int32) error
	DeleteFields(field string) error
	GetByEvent(logname string, eventid int32) ([]byte,error)
}
