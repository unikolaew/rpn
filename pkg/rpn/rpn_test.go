package rpn_test

import (
	"testing"

	"github.com/unikolaew/rpn/pkg/rpn"
)

func TestCalc(t *testing.T) {
	// Определяем структуру для тестовых случаев
	type test struct {
		name           string  // Имя теста
		expression     string  // Выражение для вычисления
		expectedResult float64 // Ожидаемый результат
		expectError    bool    // Ожидается ли ошибка
	}

	// Определяем тестовые случаи
	testCases := []test{
		{
			name:           "simple",
			expression:     "1+1", // Просто сложение
			expectedResult: 2,
			expectError:    false,
		},
		{
			name:           "priority_1",
			expression:     "(2+3)*4-5/5", // Приоритет (со скобками)
			expectedResult: 19,
			expectError:    false,
		},
		{
			name:           "priority_2",
			expression:     "2+2*2", // Приоритет (без скобок)
			expectedResult: 6,
			expectError:    false,
		},
		{
			name:           "dividing",
			expression:     "1/2", // Деление
			expectedResult: 0.5,
			expectError:    false,
		},
		{
			name:           "division_by_zero",
			expression:     "1/0", // Деление на ноль
			expectedResult: 0,     // Ожидаемый результат не важен, ожидается ошибка
			expectError:    true,
		},
		{
			name:           "invalid_expression",
			expression:     "(1+2", // Некорректное выражение (недостаток закрывающей скобки)
			expectedResult: 0,
			expectError:    true,
		},
		{
			name:           "empty_expression",
			expression:     "", // Пустое выражение
			expectedResult: 0,
			expectError:    true,
		},
		{
			name:           "unmatched_parentheses",
			expression:     "(3+2))", // Некорректное выражение (избыточная закрывающая скобка)
			expectedResult: 0,
			expectError:    true,
		},
		{
			name:           "invalid_characters",
			expression:     "1 + a * 2", // Выражение с недопустимым символом
			expectedResult: 0,
			expectError:    true,
		},
	}

	// Запускаем тесты
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := rpn.Calc(testCase.expression) // Вызываем функцию вычисления

			// Проверяем, ожидается ли ошибка
			if testCase.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none for expression: %s", testCase.expression)
				}
				return // Если ошибка ожидается, переходим к следующему тесту
			}

			// Если ошибка не ожидается, проверяем, что она не произошла
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", testCase.expression, err)
			}

			// Проверяем, что результат соответствует ожидаемому
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}
}
