package main

import (
	"github.com/asiafrolova/Final_task/orkestrator_service/internal/application"
)

func main() {
	app := application.New()
	app.RunServer()
}
