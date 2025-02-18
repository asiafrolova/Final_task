package application

import (
	"fmt"
	"net/http"
	"os"

	handlers "github.com/asiafrolova/Final_task/orkestrator_service/internal/handlers"
	logger "github.com/asiafrolova/Final_task/orkestrator_service/internal/logger"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	logger.Init()
	return &Application{config: ConfigFromEnv()}

}

func (a *Application) RunServer() error {

	logger.Info(fmt.Sprintf("Server started address: %s", a.config.Addr))
	http.HandleFunc("/api/v1/calculate", handlers.AddExpressionsHandler)
	http.HandleFunc("/api/v1/expressions", handlers.GetExpressionsListHandler)
	http.HandleFunc("/api/v1/expressions/:id", handlers.GetExpressionByIDHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
