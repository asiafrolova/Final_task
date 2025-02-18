package application

import (
	"fmt"
	"net/http"
	"os"
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
	return &Application{config: ConfigFromEnv()}

}
func AddExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add expressions")
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", AddExpressionsHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}
