package parsefield

import (
	"testing"
	"net/http"
	"encoding/json"
	"net/http/httptest"
	"github.com/julienschmidt/httprouter"
	"bytes"
	"sync"
	"io/ioutil"
)

type Body struct {
	Proc string `json:"proc"`
	Path string `json:"path"`
	Name string `json:"name"`
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
	json.NewEncoder(b).Encode(body)
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
	right := string("proc\npath\nname\n")
	if string(rbody) != right {
		t.Errorf("responce is different")
	}

}