package parsefield

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/MonaxGT/parsefields/storage"
	"github.com/julienschmidt/httprouter"
)

const (
	defaultDBName     = "parse"
	defaultTableField = "fields"
	defaultTableEvent = "events"
)

type Config struct {
	Fields    *sync.Map
	Events    *sync.Map
	separator string
	DB        storage.Database
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

func (c *Config) eventsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var events []string
	c.Events.Range(func(key, value interface{}) bool {
		events = append(events, fmt.Sprintf("%v", key))
		return true
	})

	fmt.Fprintln(w, fmt.Sprintf(strings.Join(events, "\n")))
}

func (c *Config) RestoreDBF() error {
	t, err := c.DB.RestoreFields()
	if err != nil {
		return err
	}
	for k := range t {
		c.Fields.Store(t[k].Field, 1)
	}
	return nil
}

func (c *Config) RestoreDBE() error {
	t, err := c.DB.RestoreEvents()
	if err != nil {
		return err
	}
	for k := range t {
		str := fmt.Sprintf("%s - %s", t[k].LogName, t[k].EventID)
		c.Events.Store(str, 1)
	}
	return nil
}

func Init(separator string, dbType string, dbURL string) (*Config, error) {
	var fields sync.Map
	var events sync.Map
	var db storage.Database
	switch dbType {
	case "reindexer":
		reindexer := &storage.Reindexer{
			DBName:         defaultDBName,
			NamespaceField: defaultTableField,
			NamespaceEvent: defaultTableEvent,
		}
		url := fmt.Sprintf("%s%s", dbURL, reindexer.DBName)
		if err := reindexer.Open(url); err != nil {
			return nil, err
		}
		db = reindexer
	default:
		log.Fatalln("error: you should choose db")
	}
	return &Config{
		Fields:    &fields,
		Events:    &events,
		separator: separator,
		DB:        db,
	}, nil
}

func (c *Config) eventDropHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ps.ByName("logname") == "" && ps.ByName("eventid") == "" {
		http.Error(w, "No data provided", http.StatusBadRequest)
		return
	}
	str := fmt.Sprintf("%s - %s", ps.ByName("logname"), ps.ByName("eventid"))
	c.Fields.Delete(str)
	eventID, err := strconv.ParseUint(ps.ByName("eventid"), 10, 64)
	err = c.DB.DeleteEvents(ps.ByName("logname"), int32(eventID))
	if err != nil {
		http.Error(w, "Can't delete record", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Delete")
}

func (c *Config) fieldDropHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if ps.ByName("field") == "" {
		http.Error(w, "No data provided", http.StatusBadRequest)
		return
	}
	c.Fields.Delete(ps.ByName("field"))
	err := c.DB.DeleteFields(ps.ByName("field"))
	if err != nil {
		http.Error(w, "Can't delete record", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Delete")
}

func (c *Config) Serve(addr string) error {
	if c.DB != nil {
		err := c.RestoreDBF()
		if err != nil {
			return err
		}
		err = c.RestoreDBE()
		if err != nil {
			return err
		}
	}
	router := httprouter.New()
	router.POST("/v1/json/", c.jsonHandler)
	router.POST("/v1/mjson/", c.mjsonHandler)
	router.GET("/v1/fields/", c.fieldsHandler)
	router.GET("/v1/events/", c.eventsHandler)
	router.GET("/v1/events/:logname/:eventid", c.eventDropHandler)
	router.GET("/v1/fields/:field", c.fieldDropHandler)
	log.Printf("Listening on %s", addr)
	return http.ListenAndServe(addr, router)

}
