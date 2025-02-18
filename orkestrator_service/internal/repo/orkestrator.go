package repo

import (
	"strings"
	"unicode"
)

// Первичная проверка выражения на корректность
func CheckExpression(exp string) bool {

	if strings.Contains(exp, "**") || strings.Contains(exp, "--") || strings.Contains(exp, "++") || strings.Contains(exp, "//") {
		return false
	}
	var balance int = 0
	for _, elem := range exp {
		if !(unicode.IsDigit(elem) || elem == '(' || elem == ')' || elem == '+' || elem == '-' || elem == '*' || elem == '/') {
			return false
		}
		if elem == '(' {
			balance++
		} else if elem == ')' {
			balance--
		}
		if balance < 0 {
			return false
		}
	}
	return balance == 0

}
