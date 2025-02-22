# Calculate Service

Service to parse and calculate simple arithmetic expressions.

## Prerequisites
- Installed Go version 1.23 or higher
- curl, Postman, or any similar app to work with HTTP API   

## Installation
1. Clone the repository to your machine ```git clone https://github.com/asiafrolova/Final_task.git``
2. Navigate into the project's directory ```cd Final_task/orkestrator_service```
3. Download all dependencies ```go mod tidy```
4. Run app ```go run ./cmd/main.go```
5. Repeat everything for the agent 
```
cd ../
cd Final_task/agent_service
go run ./cmd/main.go
go mod tidy
```
## How to use
**Possibilities**
- addition `+`, subtract `-`, multiplicity `*`, divide `/` operations
- any complex nested parentheses with `(` and `)`
- int and float numbers (I hope within the range -1e308..1e308) with `.` as decimal separator ()
- unary minus `-` (regular minus sign) for numbers and parentheses group's
**Examples**
*POST expression*
```
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "(3.14*2)"
}'
```
-Successfully
`{
    "id": "id1"
}
`
*GET expression by id*
```
curl --location 'localhost:8080/api/v1/expressions/id1'
```
-Successfully
`{
    "expression":
        {
            "id": "id1",
            "status": "Completed",
            "result": 6.28
        }
}`
*GET list expressions*
```
curl --location 'localhost:8080/api/v1/expressions'
```
-Successfully
`{
    "expressions": [
        {
            "id": "id1",
            "status": "Completed",
            "result": 6.28
        },
        {
            "id": "id2",
            "status": "Todo",
            "result": 0
        }
    ]
}
`
**Examples errors**
*Invalid expression*
`curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "((32*2)"
}'`

`curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "3a-1"
}'`
`curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "3++1"
}'`

*Failed status*
`curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1/0"
}'`

`curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "3.0.0-1"
}'`


   
