package repo

import (
	"fmt"
)

// Статусы выражений
var (
	TODO      string = "Todo"      //выражение ожидает
	PENDING   string = "Pending"   //выражение выполняется
	FAILED    string = "failed"    //произошла ошибка
	COMPLETED string = "completed" //выражение успешно вычислено
)

// Структура выражения
type Expression struct {
	Id     string `json:"id"`
	exp    string `json:"-"`
	Status string `json:"status"`
	Result int    `json:"result"`
}

var (
	expressionsData map[string]Expression     //Хэш таблица с выражениями (ключ - id)
	lastID          int                   = 0 //Переменная для присвоения выражениям id
)

func Init() {
	//Инициализация осуществляемая только один раз
	if expressionsData == nil {
		expressionsData = make(map[string]Expression)
	}
}

// Добавления выражения
func AddExpression(exp string) (string, error) {
	var validExp bool = CheckExpression(exp)
	if !validExp {
		return "", ErrInvalidExpression
	}
	currentID := GenerateID()
	expressionsData[currentID] = Expression{Id: currentID, exp: exp, Status: TODO}
	return currentID, nil

}

// Возврат выражения по id
func GetExpressionByID(id string) (Expression, error) {
	exp, ok := expressionsData[id]
	if ok {
		return exp, nil
	}
	return Expression{}, ErrKeyExists
}

// Возврат всех выражений
func GetExpressionsList() []Expression {
	data := make([]Expression, 0)
	for _, v := range expressionsData {
		data = append(data, v)
	}
	return data

}

// Генерация id
func GenerateID() string {
	lastID++
	return fmt.Sprintf("id%d", lastID)
}
