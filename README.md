```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)+1+9*(1+2+3)"
}'

```
```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(2+2)"
}'

```

```
curl --location 'localhost:8080/api/v1/expressions'
```
```
curl --location 'localhost:8080/api/v1/expressions?id=id1'
```

```
curl --location 'localhost:8080/internal/task'
```

```
curl --location 'localhost:8080/internal/task' \
--header 'Content-Type: application/json' \
--data '{
  "id": "se6",
  "result": 59
}'
```