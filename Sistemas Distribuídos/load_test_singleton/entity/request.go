package entity

import "fmt"

type Request struct {
	Num1      float64 `json:"num1"`
	Num2      float64 `json:"num2"`
	Operation string  `json:"operation"`
}

func (r Request) GetInput() string {
	return ""
}

func GetResponse(calculator Calculator, req Request) Response {
	value, err := calculator.Calculate(req.Num1, req.Num2, req.Operation)
	if err != nil {
		return Response{
			Result: err.Error(),
		}
	}
	return Response{
		Result: fmt.Sprintf("%s = %f", req.GetInput(), value),
	}
}
