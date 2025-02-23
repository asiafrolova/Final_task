package repo

import (
	"fmt"
	"strconv"

	orkestrator "github.com/asiafrolova/Final_task/orkestrator_service/pkg/orkestrator"
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
	expressionsData[currentID] = &orkestrator.Expression{Id: currentID, Exp: exp, Status: orkestrator.TODO}
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
	if currentExpression == nil || currentExpression.Status != orkestrator.PENDING {
		err := SetCurrentExpression()
		if err != nil {
			return
		}
	}

	for ind, elem := range currentExpression.SimpleExpressions {
		//Это выражение уже было отправлено агенту
		if elem.Processed {
			continue
		}

		newSimpleExpression, err := currentExpression.ConvertExpression(elem.Id)
		//Для выражения посчитаны нужные премененные
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

		if elem.Status == orkestrator.TODO {
			expressionsData[ind].Status = orkestrator.PENDING

			currentExpression = expressionsData[ind]
			tokenizeString, err := currentExpression.TokenizeString()
			if err != nil {

				currentExpression.Status = orkestrator.FAILED
				continue
			} else if len(tokenizeString) == 1 {
				currentExpression.Status = orkestrator.COMPLETED
				result, err := strconv.ParseFloat(tokenizeString[0], 64)
				if err != nil {
					currentExpression.Status = orkestrator.FAILED
				} else {
					currentExpression.Result = result
				}
				continue
			}
			_, err, _ = currentExpression.SplitExpression(tokenizeString)
			if err != nil {
				currentExpression.Status = orkestrator.FAILED
				continue
			}
			expressionsData[ind].WaitResult()
			return nil
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
		currentExpression.Status = orkestrator.FAILED
		SetCurrentExpression()
		return nil
	}
	err = currentExpression.SetResultSimpleExpression(id, result)
	if err != nil {
		return err
	}
	//Для выражения посчитаны все подзадачи
	if len(currentExpression.SimpleExpressions) == len(currentExpression.SimpleExpressionsResults) {
		currentExpression.Status = orkestrator.COMPLETED
		//Ответ лежит в последней подзадаче (последнем действии)
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
