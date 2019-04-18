package main

import (
	"flag"
	"log"
	"os"

	"github.com/MonaxGT/parsefields"
	"github.com/peterbourgon/ff"
)

//str := `{"process_name": "calc.exe", "process_path":"C:\\windows\\system32"}`

func main() {
	fs := flag.NewFlagSet("parsefield", flag.ExitOnError)
	var (
		listenAddr = fs.String("listen", ":8000", "Listen address")
		separator  = fs.String("sep", " -> ", "Separator nested")
		dbType     = fs.String("db-type", "reindexer", "Storage type")
		dbURL      = fs.String("db-url", "cproto://127.0.0.1:6534/", "Storage URL")
	)

	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("PARSEFIELD"))

	c, err := parsefield.Init(*separator, *dbType, *dbURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(c.Serve(*listenAddr))
}
