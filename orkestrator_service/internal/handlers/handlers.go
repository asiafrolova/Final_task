package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	logger "github.com/asiafrolova/Final_task/orkestrator_service/internal/logger"
	repo "github.com/asiafrolova/Final_task/orkestrator_service/internal/repo"
	orkestrator "github.com/asiafrolova/Final_task/orkestrator_service/pkg/orkestrator"
)

var (
	WAITING_TIME = time.Millisecond * 100 //Ожидание появление новой задачи для отправки агенту
)

// Запрос на добавление выражения
type RequestAddExpression struct {
	Expression string `json:"expression"`
}

// Ответ после успешного добавление выражения
type ResponseAddExpression struct {
	Id string `json:"id"`
}

// Ответ на запрос листа выражений
type ResponseListExpressions struct {
	Expressions []orkestrator.Expression `json:"expressions"`
}

// Ответ на запрос выражения по id
type ResponseIDExpression struct {
	Expression orkestrator.Expression `json:"expression"`
}

// Ответ на запрос задачи (простого выражения с одним действием)
type ResponseGetTask struct {
	Task orkestrator.SimpleExpression `json:"task"`
}

// Запрос на добавление результата простого выражения
type RequestAddResult struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
	Err    string  `json:"error"`
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
		if err == orkestrator.ErrInvalidExpression {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity) //Выражение не прошло первичную проверку
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError) //Неизвестная ошибка
			return
		}
	}

	response := ResponseAddExpression{Id: id}
	w.WriteHeader(http.StatusOK)
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

// Хендлер для получения списка выражений
func GetExpressionsListHandler(w http.ResponseWriter, r *http.Request) {
	//Переключение на хендлер по id, если он задан
	if r.URL.Query().Has("id") {
		GetExpressionByIDHandler(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	repo.Init()
	response := ResponseListExpressions{Expressions: repo.GetExpressionsList()}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Ошибка при сериализации
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	logger.Info("List of expressions returned successfully")

}

// Хендлер для получения выражения по id
func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	repo.Init()
	id := r.URL.Query().Get("id")
	exp, err := repo.GetExpressionByID(id)
	if err == orkestrator.ErrKeyExists {
		http.Error(w, err.Error(), http.StatusNotFound) //Нет такого ключа
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Неизвестная ошибка
		return
	}

	response := ResponseIDExpression{}
	response.Expression = *exp

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) //Ошибка при сериализации
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	logger.Info("Expression by id returned successfully")
}

// Хендлер для отправки задач агенту
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		//Переключение на хэндлер получения результата
		GetResultTaskHandle(w, r)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	timer := time.NewTimer(WAITING_TIME)
	repo.Init()
	repo.GetSimpleOperations()
	select {
	case <-timer.C:
		//Если время истекло отвечаем, что нет задач
		http.Error(w, orkestrator.ErrNotExpression.Error(), http.StatusNotFound)
		return
	case sExp := <-repo.SimpleExpressions:
		//Задача нашлась
		response := ResponseGetTask{Task: sExp}
		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) //Ошибка при сериализации

		}

		w.WriteHeader(http.StatusOK)
		w.Write(res)
		logger.Info("Task sent successfully")
	}
}

// Хендлер для получения результатов от агента
func GetResultTaskHandle(w http.ResponseWriter, r *http.Request) {
	request := RequestAddResult{}
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&request)
	repo.Init()

	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusUnprocessableEntity) //Ошибка при десериализации
		return
	}

	if request.Err != "" {
		//Пришла ошибка при вычислении
		err = repo.SetResult(request.Id, 0, fmt.Errorf(request.Err))
		if err == orkestrator.ErrKeyExists {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			logger.Info("Got an error in the task")
		}

	} else {

		err = repo.SetResult(request.Id, request.Result, nil)
		if err == orkestrator.ErrKeyExists {
			//Пришел несуществующий ключ
			http.Error(w, err.Error(), http.StatusNotFound)

		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

		} else {
			logger.Info("The task result was successfully written")
		}

	}

}
