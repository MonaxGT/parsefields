[![Codacy Badge](https://api.codacy.com/project/badge/Grade/dc9f7432a97a4d5296045487e1db44f1)](https://app.codacy.com/app/MonaxGT/parsefields?utm_source=github.com&utm_medium=referral&utm_content=MonaxGT/parsefields&utm_campaign=Badge_Grade_Dashboard)
[![GoDoc](https://godoc.org/github.com/MonaxGT/parsefields?status.png)](https://godoc.org/github.com/MonaxGT/parsefields)
[![Build Status](https://travis-ci.com/MonaxGT/gomalshare.svg?branch=master)](https://travis-ci.com/MonaxGT/gomalshare)
[![Go Report Card](https://goreportcard.com/badge/github.com/MonaxGT/parsefields)](https://goreportcard.com/report/github.com/MonaxGT/parsefields)

# parsefields

Tool for parse JSON-like logs for collecting unique fields. Main purpose to collect JSON-data with typical events and fields it is useful when you want to create mapping schema for database and you want to reduce the risks of forgotten fields. 
By default separator between to nested structs is "->", but you can change it with environment.

API consist:

```
-  /v1/json/
-  /v1/mjson/
-  /v1/fileds/
-  /v1/events/
-  /v1/events/:logname/:eventid
-  /v1/fields/:field
```

P.S. additionally info about all new events/fields will be show in stdout.

## Deploy

```sh
docker build . -t parsefield
docker run -d -p 8000:8000 parsefield
```

or 

```sh
docker-compose -p 8000:8000 -d up 
```

## Usage

### Push new log for parse

Single message per request

```sh
curl -X POST -d '{"process_name": "calc.exe", "process_path":"C:\\windows\\system32"}'  127.0.0.1:8000/v1/json/
```

Multiple message per request 

```sh
curl -X POST -d '[{"process_name": "calc.exe", "process_path":"C:\\windows\\system32"},{"process_image": "calc.exe", "process_path":"C:\\windows\\system32"},{"pid":"1"}]'  127.0.0.1:8000/v1/mjson/
```

### All unique fields

```sh
curl 127.0.0.1:8000/v1/fields/
```

### All unique events

```sh
curl 127.0.0.1:8000/v1/events/
```

### Delete events,fields

```sh
curl 127.0.0.1:8000/v1/events/Sysmon/999 - delete events with logname Sysmon and eventId 999
curl 127.0.0.1:8000/v1/fields/key - delete field with name key
```
