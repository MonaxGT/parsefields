package parsefield

import (
	"sync"
	"testing"
)

func TestParse(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := &Config{
		Fields: &fields,
		Events: &events,
	}
	body := []byte(`{"event_id":1,"event_log_name":"sysmon"}`)
	err := c.parse(body)
	if err != nil {
		t.Error(err)
	}
	body = []byte(`"event_id":1,"event_log_name":"sysmon"`)
	err = c.parse(body)
	if err == nil {
		t.Error(err)
	}

}

func TestParseMulti (t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := &Config{
		Fields: &fields,
		Events: &events,
	}
	body := []byte(`[{"event_id":1,"event_log_name":"sysmon"},{"event_id":2,"event_log_name":"security"}]`)
	err := c.parseMulti(body)
	if err != nil {
		t.Error(err)
	}
	body = []byte(`^"event_id":1,"event_log_name":"sysmon"`)
	err = c.parse(body)
	if err == nil {
		t.Error(err)
	}
}



