package rpn

import (
	"errors"
	"strconv"
	"unicode"
)

func isOperator(char rune) bool {
	operators := []rune{'+', '-', '*', '/'}
	for _, operator := range operators {
		if char == operator {
			return true
		}
	}
	return false
}

func addNumberToSlice(currentNum string, nums *[]float64) error {
	num, err := strconv.ParseFloat(currentNum, 64)
	if err != nil {
		return err
	}
	*nums = append(*nums, num)
	return nil
}

func evaluateExpression(nums []float64, operators []rune) (float64, error) {
	if len(operators) == 0 {
		return nums[0], nil
	}

	var result []float64 = append(nums, 0)

	for i := 0; i < len(operators); i++ {
		switch operators[i] {
		case '*':
			temp := result
			result = append(result[:i], result[i]*result[i+1])
			result = append(result, temp[i+2:]...)
		case '/':
			if nums[i+1] == 0 {
				return 0, errors.New("Деление на 0")
			}
			temp := result
			result = append(result[:i], result[i]/result[i+1])
			result = append(result, temp[i+2:]...)
		}
	}

	for i := 0; i < len(operators); i++ {
		switch operators[i] {
		case '+':
			temp := result
			result = append(result[:i], result[i]+result[i+1])
			result = append(result, temp[i+2:]...)
		case '-':
			temp := result
			result = append(result[:i], result[i]-result[i+1])
			result = append(result, temp[i+2:]...)
		}
	}

	return result[0], nil
}

func Calc(expression string) (float64, error) {
	if len(expression) == 0 {
		return 0, errors.New("Введено пустое выражение")
	}

	var nums []float64
	var operators []rune
	var currentNum string
	i := 0
	for i < len(expression) {
		char := rune(expression[i])

		if i == 0 && char == '-' {
			currentNum += string(char)
			continue
		}

		if unicode.IsDigit(char) || char == '.' {
			currentNum += string(char)

		} else if char == '(' {
			subExpression := ""
			j := i + 1
			openBrackets := 1
			for j < len(expression) {
				if expression[j] == '(' {
					openBrackets++
				}
				if expression[j] == ')' {
					openBrackets--
				}
				if openBrackets == 0 {
					break
				}
				subExpression += string(expression[j])
				j++
			}
			if openBrackets != 0 {
				return 0, errors.New("Есть незакрытые скобки")
			}

			result, err := Calc(subExpression)
			if err != nil {
				return 0, err
			}
			nums = append(nums, result)
			i = j

		} else if isOperator(char) {
			if currentNum != "" {
				err := addNumberToSlice(currentNum, &nums)
				if err != nil {
					return 0, err
				}
				currentNum = ""
			}
			operators = append(operators, char)

		} else {
			return 0, errors.New("Введено некорректное выражение")
		}
		i++
	}

	if currentNum != "" {
		err := addNumberToSlice(currentNum, &nums)
		if err != nil {
			return 0, err
		}
	}

	if len(nums) == 0 || len(nums) == len(operators) {
		return 0, errors.New("Введено некорректное выражение")
	}

	result, err := evaluateExpression(nums, operators)
	if err != nil {
		return 0, err
	}

	return result, nil
}
