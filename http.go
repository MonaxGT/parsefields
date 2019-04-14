package parsefield

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/julienschmidt/httprouter"
)

type Config struct {
	Fields *sync.Map
	separator string
}

func (c *Config) jsonHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read the request body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "No data provided", http.StatusBadRequest)
		return
	}
	c.parse(body)
}

func (c *Config) mjsonHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read the request body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "No data provided", http.StatusBadRequest)
		return
	}
	c.parseMulti(body)
}

func (c *Config) fieldsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var fields []string
	c.Fields.Range(func(key, value interface{}) bool {
		fields = append(fields, fmt.Sprintf("%v", key))
		return true
	})

	fmt.Fprintln(w, fmt.Sprintf(strings.Join(fields, "\n")))
}

func Init(separator string) *Config {
	var fields sync.Map
	return &Config{
		Fields: &fields,
		separator: separator,
	}
}

func (c *Config) Serve(addr string) error {
	router := httprouter.New()
	router.POST("/v1/json/", c.jsonHandler)
	router.POST("/v1/mjson/", c.mjsonHandler)
	router.GET("/v1/fields/", c.fieldsHandler)
	log.Printf("Listening on %s", addr)
	return http.ListenAndServe(addr, router)

}
