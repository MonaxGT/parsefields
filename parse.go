package parsefield

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Config) check (key string, value interface{}) {
	//fmt.Println(key, value)
	if _, ok := c.Fields.Load(key); !ok {
		c.Fields.Store(key, 1)
		fmt.Printf("New field: %s \n", key)
	}
}

func (c *Config) deep (b map[string]interface{}, prefix string) {
	for key, value := range b {
		if b, ok := value.(map[string]interface{}); ok {
			c.deep(b,fmt.Sprintf("%s%s%s",prefix,c.separator,key))
			continue
		}
		c.check(fmt.Sprintf("%s%s%s",prefix,c.separator,key),value)
	}
}

func (c *Config) parse (body []byte) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.Decode(&data)
	for key, value := range data {
		if b, ok := value.(map[string]interface{}); ok {
			c.deep(b,key)
			continue
		}
		c.check(key,value)
		}
	}

func (c *Config) parseMulti (body []byte) {
	data := []map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.Decode(&data)
	for i := range data {
		for key,value := range data[i] {
			if b, ok := value.(map[string]interface{}); ok {
				c.deep(b,key)
				continue
			}
			c.check(key,value)
		}
		}
	}