package calculation_test

import (
	"testing"

	"github.com/fstr52/calculator/pkg/calculation"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name       string
		expression string
		expResult  float64
	}{
		{
			name:       "Simple negative",
			expression: "(-5)+3",
			expResult:  -2,
		},
		{
			name:       "Simple negative",
			expression: "-5+3",
			expResult:  -2,
		},
		{
			name:       "Negative in brackets",
			expression: "5+(-3)",
			expResult:  2,
		},
		{
			name:       "Multiply negatives",
			expression: "(-5)*(-3)",
			expResult:  15,
		},
		{
			name:       "Divide negatives",
			expression: "(-5)/(-2)",
			expResult:  2.5,
		},
		{
			name:       "Sum of negative decimals",
			expression: "(-5.5)+(-4.5)",
			expResult:  -10,
		},
		{
			name:       "Complex brackets 1",
			expression: "((-5)+3)*(2+(-1))",
			expResult:  -2,
		},
		{
			name:       "Complex brackets 2",
			expression: "((2+3)*4)/((-5)+1)",
			expResult:  -5,
		},
		{
			name:       "Operations order",
			expression: "2+3*4-1",
			expResult:  13,
		},
		{
			name:       "Large nums",
			expression: "999999999999+1",
			expResult:  1000000000000,
		},
		{
			name:       "Multiple brackets",
			expression: "((((1+2)+3)+4)+5)",
			expResult:  15,
		},
		{
			name:       "Complex decimals",
			expression: "(-2.5)*(-4)*(-0.5)",
			expResult:  -5,
		},
		{
			name:       "Expression with spaces",
			expression: "  5  +  3  ",
			expResult:  8,
		},
		{
			name:       "Long",
			expression: "1+2+3+4+5+6+7+8+9+10",
			expResult:  55,
		},
		{
			name:       "Priority",
			expression: "2*3+4*5",
			expResult:  26,
		},
		{
			name:       "Decimal negative",
			expression: "-.5",
			expResult:  -0.5,
		},
		{
			name:       "Fraction",
			expression: "((1+2)*(3+4))/(5+(-2))",
			expResult:  7,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := calculation.Calculate(testCase.expression)
			if err != nil {
				t.Fatalf("got error: %s, expected nil", err)
			}
			if res != testCase.expResult {
				t.Fatalf("got: %f, expected: %f", res, testCase.expResult)
			}
		})
	}

	ErrorCases := []struct {
		name       string
		expression string
		expResult  float64
		expError   error
	}{
		{
			name:       "Division by zero",
			expression: "5/0",
			expResult:  0,
			expError:   calculation.ErrDivisionByZero,
		},
		{
			name:       "Invalid brackets 1",
			expression: "((2+3)",
			expResult:  0,
			expError:   calculation.ErrInvalidExpression,
		},
		{
			name:       "Invalid brackets 2",
			expression: "(2+3))",
			expResult:  0,
			expError:   calculation.ErrInvalidExpression,
		},
		{
			name:       "Empty expression",
			expression: "",
			expResult:  0,
			expError:   calculation.ErrInvalidExpression,
		},
	}

	for _, testCase := range ErrorCases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := calculation.Calculate(testCase.expression)
			if err == nil || err != testCase.expError {
				t.Fatalf("expected error: %s, got: %s", testCase.expError, err)
			}
			if res != testCase.expResult {
				t.Fatalf("expected: %f, but got: %f", testCase.expResult, res)
			}
		})
	}
}
