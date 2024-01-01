# Search

---

## Getting started

### Installation

1. To export environmental variables, environmental variables are stored in `.env` file, run
   ```
   export $(grep -v '^#' ./.env | xargs -d '\n')
   ```

2. To build, run 
   ```
   make build
   ```

3. To run the project on `port` and `ip` that are sets in `.env` file, run 
   ```
   ./bin/search
   ```

---

## APIs

### Search

Request:
```shell
curl -X GET http://127.0.0.1:8080/search?search=title&page_number=1&page_size=10
```

Response:
```json
{
  "books": [
    {
      "title": "title1",
      "type": "textbook",
      "cover": "cover1"
    },
    {
      "title": "title2",
      "type": "textbook",
      "cover": "cover2"
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 10
  }
}
```
or
```json
{
  "error": {
    "code": 500,
    "detail": "Internal Server Error",
    "local_detail": "خطای سرور",
    "trace_id": "01HK1BAYHV86KM3QKD91G5J2V9",
    "path": "error path",
    "trace": "error trace"
  }
}
```
