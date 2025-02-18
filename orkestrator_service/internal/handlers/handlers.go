package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	logger "github.com/asiafrolova/Final_task/orkestrator_service/internal/logger"
	repo "github.com/asiafrolova/Final_task/orkestrator_service/internal/repo"
)

// Запрос на добавление выражения
type RequestAddExpression struct {
	Expression string `json:"expression"`
}

// Ответ после успешного добавление выражения
type ResponseAddExpression struct {
	Id string `json:"id"`
}

// Запрос листа выражений
type ResponseListExpressions struct {
	Expressions []repo.Expression `json:"expressions"`
}

// Хендлер для добавления выражения
func AddExpressionsHandler(w http.ResponseWriter, r *http.Request) {

	request := new(RequestAddExpression)
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError) //Ошибка при десериализации
		return
	}
	repo.Init()
	id, err := repo.AddExpression(request.Expression)
	if err != nil {
		if err == repo.ErrInvalidExpression {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity) //Выражение не прошло первичную проверку
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError) //Неизвестная ошибка
			return
		}
	}

	response := ResponseAddExpression{Id: id}
	w.WriteHeader(201)
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Ошибка при сериализации
		return
	}
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Info(fmt.Sprintf("Expression with id %s added", response.Id))

}
func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {

}

// Хендлер для получения списка выражений
func GetExpressionsListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	repo.Init()
	response := ResponseListExpressions{Expressions: repo.GetExpressionsList()}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Ошибка при сериализации
	}
	w.WriteHeader(200)
	logger.Info(string(res))
	w.Write(res)
	logger.Info("List of expressions returned successfully")

}
