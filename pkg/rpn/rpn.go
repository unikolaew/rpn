package rpn

import (
	"errors"
	"strconv"
	"unicode"
)

// Определяем ошибки, которые могут быть возвращены в ходе вычислений.
var (
	ErrInvalidExpression = errors.New("invalid expression") // Ошибка для некорректного выражения.
	ErrDivisionByZero    = errors.New("division by zero")   // Ошибка деления на ноль.
)

// Проверяет, является ли символ оператором (+, -, *, /).
func isOperator(char rune) bool {
	operators := []rune{'+', '-', '*', '/'} // Список операторов.
	for _, operator := range operators {
		if char == operator { // Если символ является оператором, возвращаем true.
			return true
		}
	}
	return false // Если нет, возвращаем false.
}

// Добавляет число в срез чисел и обрабатывает ошибки парсинга.
func addNumberToSlice(currentNum string, nums *[]float64) error {
	num, err := strconv.ParseFloat(currentNum, 64) // Парсит строку в число с плавающей запятой.
	if err != nil {
		return err // Возвращает ошибку, если парсинг не удался.
	}
	*nums = append(*nums, num) // Добавляем число в срез.
	return nil
}

// Оценивает выражение с использованием чисел и операторов.
func evaluateExpression(nums []float64, operators []rune) (float64, error) {
	var results []float64 // Срез для хранения промежуточных результатов
	var ops []rune        // Срез для хранения операторов сложения и вычитания

	// Начинаем с первого числа
	results = append(results, nums[0])
	for i := 0; i < len(operators); i++ {
		// Проверяем, является ли текущий оператор умножением или делением
		if operators[i] == '*' || operators[i] == '/' {
			if operators[i] == '*' {
				// Выполняем умножение
				last := results[len(results)-1]                       // Последний результат
				next := nums[i+1]                                     // Следующее число
				results = append(results[:len(results)-1], last*next) // Обновляем последний результат
			} else if operators[i] == '/' {
				// Проверяем деление на ноль
				if nums[i+1] == 0 {
					return 0, ErrDivisionByZero // Возвращаем ошибку
				}
				last := results[len(results)-1]                       // Последний результат
				next := nums[i+1]                                     // Следующее число
				results = append(results[:len(results)-1], last/next) // Обновляем последний результат
			}
		} else {
			// Если оператор сложения или вычитания, добавляем его в срез ops
			ops = append(ops, operators[i])
			results = append(results, nums[i+1]) // Добавляем следующее число в результаты
		}
	}

	// Обрабатываем сложение и вычитание
	result := results[0] // Начинаем с первого результата
	for i, op := range ops {
		if op == '+' {
			result += results[i+1] // Сложение
		} else if op == '-' {
			result -= results[i+1] // Вычитание
		}
	}

	return result, nil // Возвращаем итоговый результат
}

// Основная функция для вычисления выражения.
func Calc(expression string) (float64, error) {
	if len(expression) == 0 { // Проверка на пустое выражение.
		return 0, ErrInvalidExpression
	}

	var nums []float64    // Срез для чисел.
	var operators []rune  // Срез для операторов.
	var currentNum string // Текущая строка числа.
	i := 0
	// Обработка каждого символа в выражении.
	for i < len(expression) {
		char := rune(expression[i]) // Приводим символ к типу rune.

		if i == 0 && char == '-' { // Обработка отрицательных чисел.
			currentNum += string(char) // Добавляем символ '-' к текущему числу.
			continue
		}

		if unicode.IsDigit(char) || char == '.' { // Проверка на цифровые символы или точку.
			currentNum += string(char) // Добавление символа к текущему числу.

		} else if char == '(' { // Обработка подвыражений в скобках.
			subExpression := "" // Инициализация подвыражения.
			j := i + 1
			openBrackets := 1 // Счетчик открывающих скобок.
			for j < len(expression) {
				if expression[j] == '(' {
					openBrackets++ // Увеличиваем счетчик для открывающей скобки.
				}
				if expression[j] == ')' {
					openBrackets-- // Уменьшаем счетчик для закрывающей скобки.
				}
				if openBrackets == 0 {
					break // Если скобки уравновешены, выходим из цикла.
				}
				subExpression += string(expression[j]) // Добавляем символ к подвыражению.
				j++
			}
			if openBrackets != 0 { // Проверка на некорректно расставленные скобки.
				return 0, ErrInvalidExpression
			}

			// Рекурсивный вызов для вычисления подвыражения.
			result, err := Calc(subExpression)
			if err != nil {
				return 0, err // Возвращаем ошибку, если возникла проблема.
			}
			nums = append(nums, result) // Добавляем результат подвыражения в список чисел.
			i = j                       // Обновляем индекс для продолжения обработки выражения.

		} else if isOperator(char) { // Если символ является оператором.
			if currentNum != "" { // Если есть текущее число, добавляем его в срез.
				err := addNumberToSlice(currentNum, &nums)
				if err != nil {
					return 0, err // Обрабатываем ошибки при добавлении числа.
				}
				currentNum = "" // Сбрасываем текущую строку числа.
			}
			operators = append(operators, char) // Добавляем оператор в срез.

		} else {
			return 0, ErrInvalidExpression // Возвращаем ошибку, если символ некорректен.
		}
		i++
	}

	// Обработка текущего числа в конце цикла.
	if currentNum != "" {
		err := addNumberToSlice(currentNum, &nums)
		if err != nil {
			return 0, err // Обрабатываем ошибки.
		}
	}

	// Проверка на корректность количества чисел и операторов.
	if len(nums) == 0 || len(nums) == len(operators) {
		return 0, ErrInvalidExpression // Возвращаем ошибку, если некорректное количество.
	}

	// Вычисление итогового результата.
	result, err := evaluateExpression(nums, operators)
	if err != nil {
		return 0, err // Обрабатываем ошибки при вычислениях.
	}

	return result, nil // Возвращаем итоговый результат.
}
