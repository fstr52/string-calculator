package calculation

import "errors"

var (
	ErrInvalidExpression = errors.New("invalid expression") // Введено неверное выражение
	ErrDivisionByZero    = errors.New("division by zero")   // Попытка произвести деление на 0
)
