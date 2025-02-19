package repo

import (
	"fmt"

	orkestrator "github.com/asiafrolova/Final_task/orkestrator_service/pkg/orkestrator"
)

// Статусы выражений
var (
	TODO      string = "Todo"      //выражение ожидает
	PENDING   string = "Pending"   //выражение выполняется
	FAILED    string = "failed"    //произошла ошибка
	COMPLETED string = "completed" //выражение успешно вычислено
)

var (
	expressionsData   map[string]*orkestrator.Expression     //Хэш таблица с выражениями (ключ - id)
	lastID            int                                = 0 //Переменная для присвоения выражениям id
	SimpleExpressions chan orkestrator.SimpleExpression      //Канал для простых выражений
	currentExpression *orkestrator.Expression                //выражение с которым сейчас работаем
)

func Init() {
	//Инициализация осуществляемая только один раз
	if expressionsData == nil {
		expressionsData = make(map[string]*orkestrator.Expression)
	}
	if SimpleExpressions == nil {
		SimpleExpressions = make(chan orkestrator.SimpleExpression)
	}
}

// Добавления выражения
func AddExpression(exp string) (string, error) {
	var validExp bool = orkestrator.CheckExpression(exp)
	if !validExp {
		return "", orkestrator.ErrInvalidExpression
	}
	currentID := GenerateID()
	expressionsData[currentID] = &orkestrator.Expression{Id: currentID, Exp: exp, Status: TODO}
	return currentID, nil

}

// Возврат выражения по id
func GetExpressionByID(id string) (*orkestrator.Expression, error) {
	exp, ok := expressionsData[id]
	if ok {
		return exp, nil
	}
	return &orkestrator.Expression{}, orkestrator.ErrKeyExists
}

// Возврат всех выражений
func GetExpressionsList() []orkestrator.Expression {
	data := make([]orkestrator.Expression, 0)
	for _, v := range expressionsData {
		data = append(data, *v)
	}
	return data

}

// просим положить еще простых выражений канал
func GetSimpleOperations() {
	if currentExpression == nil || currentExpression.Status != PENDING {
		err := SetCurrentExpression()
		if err != nil {
			return
		}
	}

	for ind, elem := range currentExpression.SimpleExpressions {
		if elem.Processed {
			continue
		}

		newSimpleExpression, err := currentExpression.ConvertExpression(elem.Id)
		if err == nil {
			currentExpression.SimpleExpressions[ind].Processed = true
			go func() {
				SimpleExpressions <- newSimpleExpression
			}()
		}
	}
}

// Устанавливаем выражение с которым будем работать
func SetCurrentExpression() error {
	for ind, elem := range expressionsData {
		if elem.Status == TODO {
			expressionsData[ind].Status = PENDING

			currentExpression = expressionsData[ind]
			tokenizeString, err := currentExpression.TokenizeString()
			if err != nil {

				currentExpression.Status = FAILED
				continue
			}
			_, err, _ = currentExpression.SplitExpression(tokenizeString)

			return err
		}
	}
	return orkestrator.ErrNotExpression
}

// Устанавливаем результат вычисления простого выражения и обновляем статус родительского выражения
func SetResult(id string, result float64, err error) error {
	if currentExpression == nil {
		return orkestrator.ErrNotExpression
	}
	if err != nil {
		currentExpression.Status = FAILED
		SetCurrentExpression()
		return orkestrator.ErrInvalidExpression
	}
	err = currentExpression.SetResultSimpleExpression(id, result)
	if err != nil {
		return err
	}

	if len(currentExpression.SimpleExpressions) == len(currentExpression.SimpleExpressionsResults) {
		currentExpression.Status = COMPLETED
		currentExpression.Result = currentExpression.SimpleExpressionsResults[currentExpression.SimpleExpressions[len(currentExpression.SimpleExpressions)-1].Id]
		SetCurrentExpression()
	}
	return nil

}

// Генерация id
func GenerateID() string {
	lastID++
	return fmt.Sprintf("id%d", lastID)
}
