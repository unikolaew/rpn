package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/unikolaew/rpn/internal/handler"
)

type testCase struct {
	name             string                      // Имя теста
	expression       string                      // Выражение для тестирования
	expectedStatus   int                         // Ожидаемый статус ответа
	expectedResponse handler.CalculationResponse // Ожидаемый ответ
}

func TestHandleCalculation(t *testing.T) {
	tests := []testCase{
		{
			name:             "simple_addition",
			expression:       "2+2",
			expectedStatus:   http.StatusOK,
			expectedResponse: handler.CalculationResponse{Result: fmt.Sprintf("%f", float64(4))},
		},
		{
			name:             "simple_subtraction",
			expression:       "5-3",
			expectedStatus:   http.StatusOK,
			expectedResponse: handler.CalculationResponse{Result: fmt.Sprintf("%f", float64(2))},
		},
		{
			name:             "simple_multiplication",
			expression:       "3*4",
			expectedStatus:   http.StatusOK,
			expectedResponse: handler.CalculationResponse{Result: fmt.Sprintf("%f", float64(12))},
		},
		{
			name:             "simple_division",
			expression:       "10/2",
			expectedStatus:   http.StatusOK,
			expectedResponse: handler.CalculationResponse{Result: fmt.Sprintf("%f", float64(5))},
		},
		{
			name:             "division_by_zero",
			expression:       "1/0", // Деление на ноль
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: handler.CalculationResponse{Error: "Internal server error"},
		},
		{
			name:             "invalid_characters",
			expression:       "2+a2", // Некорректное выражение
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: handler.CalculationResponse{Error: "Expression is not valid"},
		},
		{
			name:             "missing_closing_parenthesis",
			expression:       "(1+2", // Пропущенная закрывающая скобка
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: handler.CalculationResponse{Error: "Expression is not valid"},
		},
		{
			name:             "unmatched_parentheses",
			expression:       "(3+2))", // Избыточная закрывающая скобка
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: handler.CalculationResponse{Error: "Expression is not valid"},
		},
		{
			name:             "empty_expression",
			expression:       "", // Пустое выражение
			expectedStatus:   http.StatusUnprocessableEntity,
			expectedResponse: handler.CalculationResponse{Error: "Expression is not valid"},
		},
		{
			name:             "complex_expression",
			expression:       "(2+3)*4-5/5", // Сложное корректное выражение
			expectedStatus:   http.StatusOK,
			expectedResponse: handler.CalculationResponse{Result: fmt.Sprintf("%f", float64(19))},
		},
	}

	for _, test := range tests {
		reqBody, _ := json.Marshal(handler.CalculationRequest{Expression: test.expression})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		w := httptest.NewRecorder()

		// Вызываем обработчик
		handler.HandleCalculation(w, req)

		// Проверяем код статуса ответа
		if w.Code != test.expectedStatus {
			t.Errorf("[%s] handler returned wrong status code: got %v want %v", test.name, w.Code, test.expectedStatus)
		}

		var response handler.CalculationResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Errorf("[%s] error when decoding response body: %v", test.name, err)
		}

		// Проверяем, что ответ совпадает с ожидаемым
		if response != test.expectedResponse {
			t.Errorf("[%s] handler returned unexpected body: got %#v want %#v", test.name, response, test.expectedResponse)
		}
	}
}
