```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)+1+9*(1+2+3.00)"
}'

```
```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(3.14*2)"
}'

```

```
curl --location 'localhost:8080/api/v1/expressions'
```
```
curl --location 'localhost:8080/api/v1/expressions?id=id3'

curl --location 'localhost:8080/api/v1/expressions/id3'
```

```
curl --location 'localhost:8080/internal/task'
```

```
curl --location 'localhost:8080/internal/task' \
--header 'Content-Type: application/json' \
--data '{
  "id": "se2",
  "result": 0,
  "error":"err"
 
}'
```

```
agent_service/cmd go run main.go
```
```
go run ./orkestrator_service/cmd/main.go
```