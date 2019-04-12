# parsefields
Tools for parse JSON-like logs for collecting unique fields

## Deploy:

```
docker build . -t parsefield
docker run -d -p 8000:8000 parsefield
```
## Usage:

Example:

### Push new log for parse

```
curl -X POST -d '{"test4": "calc.exe", "process_path":"C:\\windows\\system32"}'  127.0.0.1:8000/v1/json/
```

### Check unique fields

```
curl 127.0.0.1:8000/v1/fields/
```
