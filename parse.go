package parsefield

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MonaxGT/parsefields/storage"
	"log"
)

const eventID = "event_id"
const logName = "event_log_name"

func (c *Config) check(key string, value interface{}) {
	//fmt.Println(key, value)
	if _, ok := c.Fields.Load(key); !ok {
		c.Fields.Store(key, 1)
		if c.DB != nil {
			err := c.DB.InsertFields(&storage.Fields{
				Field: key,
			})

			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("New field: %s \n", key)
	}
}

func (c *Config) deep(b map[string]interface{}, prefix string) {
	for key, value := range b {
		if b, ok := value.(map[string]interface{}); ok {
			c.deep(b, fmt.Sprintf("%s%s%s", prefix, c.separator, key))
			continue
		}
		c.check(fmt.Sprintf("%s%s%s", prefix, c.separator, key), value)
	}
}

func (c *Config) parse(body []byte) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.Decode(&data)
	for key, value := range data {
		if b, ok := value.(map[string]interface{}); ok {
			c.deep(b, key)
			continue
		}
		c.check(key, value)
		if key == logName {
			c.parseEvent(data)
		}
	}
}

func (c *Config) parseEvent(data map[string]interface{}) {
	str := fmt.Sprintf("%s - %d", data[logName], int32(data[eventID].(float64)))
	if _, ok := c.Events.Load(str); !ok {
		c.Events.Store(str, 1)
		if c.DB != nil {
			err := c.DB.InsertEvents(&storage.Events{
				Data:    data,
				EventID: int32(data[eventID].(float64)),
				LogName: data[logName].(string),
			})
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Printf("New event add: %s \n", str)
	}
}

func (c *Config) parseMulti(body []byte) {
	data := []map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.Decode(&data)
	for i := range data {
		for key, value := range data[i] {
			if b, ok := value.(map[string]interface{}); ok {
				c.deep(b, key)
				continue
			}
			c.check(key, value)
			if key == logName {
				c.parseEvent(data[i])
			}
		}
	}
}
