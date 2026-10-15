package entity

import (
	"errors"
	"math"
)

type Calculator interface {
	Calculate(num1, num2 float64, operation string) (float64, error)
}

type SimpleCalculator struct {
}

func (s SimpleCalculator) Calculate(num1, num2 float64, operation string) (float64, error) {
	switch operation {
	case "+":
		return s.add(num1, num2)
	case "-":
		return s.sub(num1, num2)
	case "*":
		return s.mult(num1, num2)
	case "/":
		return s.div(num1, num2)
	default:
		return 0.0, errors.New("operation not supported")
	}
}

func (s SimpleCalculator) add(num1, num2 float64) (float64, error) {
	return num1 + num2, nil
}

func (s SimpleCalculator) sub(num1, num2 float64) (float64, error) {
	return num1 - num2, nil
}

func (s SimpleCalculator) mult(num1, num2 float64) (float64, error) {
	return num1 * num2, nil
}

func (s SimpleCalculator) div(num1, num2 float64) (float64, error) {
	if num2 == 0 {
		return 0, errors.New("division by zero")
	}
	return num1 / num2, nil
}

type ScienceCalculator struct {
}

func (s ScienceCalculator) Calculate(num1, num2 float64, operation string) (float64, error) {
	switch operation {
	case "**":
		return s.exp(num1, num2)
	default:
		return 0.0, errors.New("operation not supported")
	}
}

func (s ScienceCalculator) exp(num1 float64, num2 float64) (float64, error) {
	return math.Pow(num1, num2), nil
}

type AnotherCalculator struct {
}

func (c AnotherCalculator) Calculate(num1, num2 float64, operation string) (float64, error) {
	switch operation {
	case "@":
		return c.antExp(num1, num2)
	default:
		return 0.0, errors.New("operation not supported")
	}
}

func (c AnotherCalculator) antExp(num1 float64, num2 float64) (float64, error) {
	return math.Pow(num1+num2, 2), nil
}
