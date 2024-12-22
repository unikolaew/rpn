package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/unikolaew/rpn/pkg/rpn"
)

// Запрос
type CalculationRequest struct {
	Expression string `json:"expression"` // Выражение
}

// Ответ
type CalculationResponse struct {
	Result string `json:"result,omitempty"` // Результат вычисления (опционально)
	Error  string `json:"error,omitempty"`  // Сообщение об ошибке (опционально)
}

// Обработчик HTTP-запросов для вычисления выражений.
func HandleCalculation(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок ответа на JSON
	w.Header().Set("Content-Type", "application/json")

	// Проверяем метод запроса, должен быть POST
	if r.Method != http.MethodPost {
		response := CalculationResponse{Error: "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError) // Устанавливаем статус ответа 500
		json.NewEncoder(w).Encode(response)           // Отправляем ответ
		return
	}

	defer r.Body.Close()       // В конце закрываем тело запроса
	var req CalculationRequest // Создаем переменную для хранения входного запроса
	// Декодируем JSON из тела запроса в переменную req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := CalculationResponse{Error: "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError) // Устанавливаем статус ответа 500
		json.NewEncoder(w).Encode(response)           // Отправляем ответ
		return
	}

	// Пытаемся вычислить выражение
	result, err := rpn.Calc(req.Expression)
	// Если получена ошибка о некорректном выражении, возвращаем 422
	if err == rpn.ErrInvalidExpression {
		response := CalculationResponse{Error: "Expression is not valid"}
		w.WriteHeader(http.StatusUnprocessableEntity) // Устанавливаем статус ответа 422
		json.NewEncoder(w).Encode(response)           // Отправляем ответ
		return
	} else if err != nil {
		// Если возникла другая ошибка, возвращаем 500
		response := CalculationResponse{Error: "Internal server error"}
		w.WriteHeader(http.StatusInternalServerError) // Устанавливаем статус ответа 500
		json.NewEncoder(w).Encode(response)           // Отправляем ответ
		return
	}

	// Если расчет успешен, формируем ответ с результатом
	response := CalculationResponse{Result: fmt.Sprintf("%f", result)}
	w.WriteHeader(http.StatusOK)        // Устанавливаем статус ответа 200
	json.NewEncoder(w).Encode(response) // Отправляем ответ с результатом
}
