package parsefield

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"fmt"
	"github.com/julienschmidt/httprouter"
)

type Body struct {
	Proc    string `json:"proc"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	EventID uint32 `json:"event_id,omitempty"`
	LogName string `json:"event_log_name,omitempty"`
}

func TestJSONHandler(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := Config{
		Fields: &fields,
		Events: &events,
	}
	router := httprouter.New()
	router.POST("/v1/json/", c.JSONHandler)
	router.GET("/v1/fields/", c.FieldsHandler)

	body := Body{
		Proc: "calc.exe",
		Path: "C:/windows/",
		Name: "gopher",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1/json/", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "http://127.0.0.1:8000/v1/fields/", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var rbody []byte
	if rbody, err = ioutil.ReadAll(rr.Body); err != nil {
		if err != nil {
			t.Errorf("can't read the request body")
		}
	}
	if rbody == nil {
		t.Errorf("server return nil data")
	}
}

func TestMJSONHandler(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := Config{
		Fields: &fields,
		Events: &events,
	}
	router := httprouter.New()
	router.POST("/v1/mjson/", c.MJSONHandler)
	router.GET("/v1/fields/", c.FieldsHandler)

	body := []Body{
		{
			Proc: "calc.exe",
			Path: "C:/windows/",
			Name: "gopher",
		},
		{
			Proc: "word.exe",
			Path: "C:/windows/",
			Name: "tester",
		},
	}
	fmt.Println(body)
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1/mjson/", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("GET", "http://127.0.0.1:8000/v1/fields/", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var rbody []byte
	if rbody, err = ioutil.ReadAll(rr.Body); err != nil {
		if err != nil {
			t.Errorf("can't read the request body")
		}
	}
	if rbody == nil {
		t.Errorf("server return nil data")
	}
}

func TestEventHandler(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := Config{
		Fields: &fields,
		Events: &events,
	}
	router := httprouter.New()
	router.POST("/v1/json/", c.JSONHandler)
	router.GET("/v1/events/", c.eventsHandler)

	body := Body{
		Proc:    "calc.exe",
		Path:    "C:/windows/",
		Name:    "gopher",
		EventID: 1,
		LogName: "Sysmon",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1/json/", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	req, err = http.NewRequest("GET", "http://127.0.0.1:8000/v1/events/", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var rbody []byte
	if rbody, err = ioutil.ReadAll(rr.Body); err != nil {
		if err != nil {
			t.Errorf("can't read the request body")
		}
	}
	if rbody == nil {
		t.Errorf("server return nil data")
	}
}

func TestEventDropHandler(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := Config{
		Fields: &fields,
		Events: &events,
	}
	router := httprouter.New()
	router.POST("/v1/json/", c.JSONHandler)
	router.DELETE("/v1/events/:logname/:eventid", c.eventDropHandler)

	body := Body{
		Proc:    "calc.exe",
		Path:    "C:/windows/",
		Name:    "gopher",
		EventID: 1,
		LogName: "Sysmon",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1/json/", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	req, err = http.NewRequest("DELETE", "http://127.0.0.1:8000/v1/events/Sysmon/1", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestFieldDropHandler(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := Config{
		Fields: &fields,
		Events: &events,
	}
	router := httprouter.New()
	router.POST("/v1/json/", c.JSONHandler)
	router.DELETE("/v1/fields/:field", c.fieldDropHandler)

	body := Body{
		Proc:    "calc.exe",
		Path:    "C:/windows/",
		Name:    "gopher",
		EventID: 1,
		LogName: "Sysmon",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1/json/", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	req, err = http.NewRequest("DELETE", "http://127.0.0.1:8000/v1/fields/proc", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestEventsBodyHandler(t *testing.T) {
	var fields sync.Map
	var events sync.Map
	c := Config{
		Fields: &fields,
		Events: &events,
	}
	router := httprouter.New()
	router.POST("/v1/json/", c.JSONHandler)
	router.GET("/v1/events/:logname/:eventid", c.eventsBodyHandler)

	body := Body{
		Proc:    "calc.exe",
		Path:    "C:/windows/",
		Name:    "gopher",
		EventID: 1,
		LogName: "Sysmon",
	}
	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://127.0.0.1:8000/v1/json/", b)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	req, err = http.NewRequest("GET", "http://127.0.0.1:8000/v1/events/Sysmon/1", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if _, err = ioutil.ReadAll(rr.Body); err != nil {
		if err != nil {
			t.Errorf("can't read the request body")
		}
	}
}
