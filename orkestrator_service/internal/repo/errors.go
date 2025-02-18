package repo

import "errors"

var (
	ErrInvalidExpression = errors.New("Invaild expression") //Ошибка в выражении
	ErrKeyExists         = errors.New("Invalid id")         //Выражения с таким id не существует
)
