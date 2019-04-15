# parsefields
Tool for parse JSON-like logs for collecting unique fields. It is useful when you want to create mapping schema for database and you want to reduce the risks of forgotten fields. 

API consist:

```json
/v1/json/   - collector for single json request
/v1/mjson/  - collector for batch json request
/v1/fileds/ - Print all fields 
/v1/events/ - Print all events
/v1/events/:logname/:eventid - Drop event
/v1/fields/:field - Drop field
```

## Deploy:

```sh
docker build . -t parsefield
docker run -d -p 8000:8000 parsefield
```

or 

```sh
docker-compose -p 8000:8000 -d up 
```

## Usage:

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


By default separator between to nested structs is "->", but you can change it with environment
