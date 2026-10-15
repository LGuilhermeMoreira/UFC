package model

import "fmt"

type Request struct {
	Num1   float64
	Num2   float64
	Action string
}

type Response struct {
	Result float64
}

func (r Response) Show() {
	fmt.Printf("Result: %f\n", r.Result)
}
