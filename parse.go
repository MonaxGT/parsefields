package parsefield

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (c *Config) parse(body []byte) {
	data := map[string]interface{}{}
	dec := json.NewDecoder(bytes.NewBuffer(body))
	dec.Decode(&data)
	for key, _ := range data {
		if _, ok := c.Fields.Load(key); ok {
			continue
		}
		c.Fields.Store(key, 1)
		fmt.Printf("New field: %s \n", key)

	}
}
