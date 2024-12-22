package application_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fstr52/string-calculator/internal/application"
)

func TestRunServer(t *testing.T) {
	testCases := []struct {
		name         string
		method       string
		input        string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid POST request",
			method:       "POST",
			input:        `{"expression": "2+2"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"4"}`,
		},
		{
			name:         "Invalid GET request",
			method:       "GET",
			input:        `{"expression": "2+2"}`,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: `{"error":"Method not allowed"}`,
		},
		{
			name:         "Invalid PUT request",
			method:       "PUT",
			input:        `{"expression": "2+2"}`,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: `{"error":"Method not allowed"}`,
		},
		{
			name:         "Valid expression",
			method:       "POST",
			input:        `{"expression": "2+2"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"4"}`,
		},
		{
			name:         "Valid expression",
			method:       "POST",
			input:        `{"expression": "(5+3)*2"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"16"}`,
		},
		{
			name:         "Valid expression",
			method:       "POST",
			input:        `{"expression": "10.5/2"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"5.25"}`,
		},
		{
			name:         "Valid expression",
			method:       "POST",
			input:        `{"expression": "  5  +  5  "}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"10"}`,
		},
		{
			name:         "Invalid JSON",
			method:       "POST",
			input:        `{"expression": "2+2"`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid. Ivalind JSON"}`,
		},
		{
			name:         "Empty expression",
			method:       "POST",
			input:        `{"expression": ""}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid"}`,
		},
		{
			name:         "Division by zero",
			method:       "POST",
			input:        `{"expression": "5/0"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid"}`,
		},
		{
			name:         "Invalid expression",
			method:       "POST",
			input:        `{"expression": "2++2"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid"}`,
		},
		{
			name:         "Invalid expression",
			method:       "POST",
			input:        `{}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid"}`,
		},
		{
			name:         "Invalid expression",
			method:       "POST",
			input:        `{"expression": "2+a"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid"}`,
		},
		{
			name:         "Invalid expression",
			method:       "POST",
			input:        `{"expression": "(2+2"}`,
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: `{"error":"Expression is not valid"}`,
		},
		{
			name:         "Large numbers",
			method:       "POST",
			input:        `{"expression": "999999999+1"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"1000000000"}`,
		},
		{
			name:         "Multiple operations",
			method:       "POST",
			input:        `{"expression": "2+2*3-4/2"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"result":"6"}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testApp := application.New()
			handler := http.HandlerFunc(application.New().RequestHandler(testApp.LoggingHandler(testApp.CalcHandler)))
			req, err := http.NewRequest(testCase.method, "api/v1/calculate", strings.NewReader(testCase.input))
			if err != nil {
				t.Fatalf("failed to create new request, %s", err)
			}

			req.Header.Set("Content-Type", "application/json")
			record := httptest.NewRecorder()
			handler.ServeHTTP(record, req)

			if record.Code != testCase.expectedCode {
				t.Errorf("expected code: %v, but got: %v", testCase.expectedCode, record.Code)
			}
			if strings.TrimSpace(record.Body.String()) != testCase.expectedBody {
				t.Errorf("expected body: %v, but got: %v", testCase.expectedBody, record.Body.String())
			}
		})
	}
}
