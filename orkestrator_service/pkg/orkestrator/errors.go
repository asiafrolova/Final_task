package orkestrator

import "errors"

var (
	ErrInvalidExpression = errors.New("Invaild expression")                     //Ошибка в выражении
	ErrKeyExists         = errors.New("Invalid id")                             //Выражения с таким id не существует
	ErrNotResult         = errors.New("Expression has not yet been calculated") //Выражение ещё не посчитано
	ErrNotExpression     = errors.New("Not expressions")                        //Нет выражений для вычисления
)
