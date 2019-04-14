package main

import (
	"flag"
	"log"
	"os"

	"github.com/MonaxGT/parsefield"
	"github.com/peterbourgon/ff"
)

//str := `{"process_name": "calc.exe", "process_path":"C:\\windows\\system32"}`

func main() {
	fs := flag.NewFlagSet("parsefield", flag.ExitOnError)
	var (
		listenAddr = fs.String("listen", ":8000", "Listen address")
		separator = fs.String("sep", " -> ", "Separator nested")
	)

	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("PARSEFIELD"))

	c := parsefield.Init(*separator)
	log.Fatal(c.Serve(*listenAddr))
}
