package main

import (
	"net/http"

	"github.com/unikolaew/rpn/internal/handler"
)

func main() {
	// Устанавливаем обработчик для маршрута /api/v1/calculate
	http.HandleFunc("/api/v1/calculate", handler.HandleCalculation)
	// Запускаем HTTP-сервер на порту 8080
	http.ListenAndServe(":8080", nil)
}
