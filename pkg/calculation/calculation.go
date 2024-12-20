package calculation

import (
	"strconv"
	"strings"
	"unicode"
)

var priors = map[string]int{
	"(": 0,
	")": 0,
	"+": 0,
	"-": 0,
	"*": 1,
	"/": 1,
}

// Запись полученной на вводе строки в обратную польскую нотацию
func toPolandNotation(expression string) (string, error) {
	resultString := ""
	operStack := ""
	var currentNumber string

	openBrackets := 0

	for i := 0; i < len(expression); i++ {
		char := expression[i]

		if char == '-' && (i == 0 || expression[i-1] == '(') {
			if i < len(expression)-1 && (unicode.IsDigit(rune(expression[i+1])) || expression[i+1] == '.') {
				currentNumber = "-"
			} else {
				resultString += "0 "
				operStack = string(char) + operStack
			}
			continue
		}

		if unicode.IsDigit(rune(char)) || char == '.' {
			currentNumber += string(char)

			if i == len(expression)-1 || (!unicode.IsDigit(rune(expression[i+1])) && expression[i+1] != '.') {
				resultString += currentNumber + " "
				currentNumber = ""
			}
			continue
		}

		if char == ' ' {
			continue
		}

		switch char {
		case '(':
			operStack = string(char) + operStack
			openBrackets++
		case ')':
			for len(operStack) > 0 && operStack[0] != '(' {
				resultString += string(operStack[0]) + " "
				operStack = operStack[1:]
			}
			if len(operStack) > 0 {
				operStack = operStack[1:]
			}
			openBrackets--
		case '+', '-', '*', '/':
			for len(operStack) > 0 && operStack[0] != '(' &&
				priors[string(operStack[0])] >= priors[string(char)] {
				resultString += string(operStack[0]) + " "
				operStack = operStack[1:]
			}
			operStack = string(char) + operStack
		}
	}

	if openBrackets != 0 {
		return "", ErrInvalidExpression
	}

	for len(operStack) > 0 {
		if operStack[0] != '(' {
			resultString += string(operStack[0]) + " "
		}
		operStack = operStack[1:]
	}

	return strings.TrimSpace(resultString), nil
}

// Проведение операций в вводной строке, записанной в форме обратной польской нотации (рассчет)
func Calculate(input string) (float64, error) {
	var signCount int

	if len(input) == 0 {
		return 0, ErrInvalidExpression
	}

	toCalc, err := toPolandNotation(input)
	if err != nil {
		return 0, err
	}

	tokens := strings.Fields(toCalc)
	stack := make([]float64, 0)

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			if len(stack) < 2 {
				return 0, ErrInvalidExpression
			}
			signCount++
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				result = a / b
			}
			stack = append(stack, result)
		default:
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, ErrInvalidExpression
			}
			stack = append(stack, num)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}
