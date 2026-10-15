package handler

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sockets/model"
	"strconv"
	"strings"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
	}
	fmt.Println(string(msg[:n]), conn.RemoteAddr())
}

func HandleCalculate(conn net.Conn) {
	defer conn.Close()
	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
		handleError(conn, err)
		return
	}
	numericalExpression := string(msg[:n])
	result, err := calculate(numericalExpression)
	if err != nil {
		handleError(conn, err)
		return
	}
	response := fmt.Sprintf("%f", result)
	_, err = conn.Write([]byte(response))
	if err != nil {
		handleError(conn, err)
		return
	}
}

func HandleCalculateWithStruct(conn net.Conn) {
	defer conn.Close()
	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
		handleError(conn, err)
		return
	}
	request, err := HandleReadRequestData(msg[:n])
	if err != nil {
		handleError(conn, err)
		return
	}
	result, err := calculateWithData(*request)
	if err != nil {
		handleError(conn, err)
		return
	}
	response, err := HandleWriteResponseData(result)
	if err != nil {
		handleError(conn, err)
		return
	}
	_, err = conn.Write(response)
	if err != nil {
		handleError(conn, err)
		return
	}
}

func calculateWithData(request model.Request) (float64, error) {
	switch request.Action {
	case "+":
		return request.Num1 + request.Num2, nil
	case "-":
		return request.Num1 - request.Num2, nil
	case "*":
		return request.Num1 * request.Num2, nil
	case "/":
		return request.Num1 / request.Num2, nil
	default:
		return 0, errors.New("invalid action")
	}
}

func handleError(conn net.Conn, err error) {
	conn.Write([]byte(err.Error()))
}

func calculate(s string) (float64, error) {
	array := strings.Split(s, " ")
	if len(array) != 3 {
		return 0, errors.New("numerical expression is invalid")
	}
	switch array[1] {
	case "+":
		return handleAddition(array)
	case "-":
		return handleSubtraction(array)
	case "*":
		return handleMultiplication(array)
	case "/":
		return handleDivision(array)
	default:
		return 0.0, errors.New("numerical expression is invalid")
	}
}

func handleAddition(array []string) (float64, error) {
	num1, err := strconv.ParseFloat(array[0], 64)
	if err != nil {
		return 0.0, err
	}
	num2, err := strconv.ParseFloat(array[2], 64)
	if err != nil {
		return 0.0, err
	}
	return num1 + num2, nil
}
func handleSubtraction(array []string) (float64, error) {
	num1, err := strconv.ParseFloat(array[0], 64)
	if err != nil {
		return 0.0, err
	}
	num2, err := strconv.ParseFloat(array[2], 64)
	if err != nil {
		return 0.0, err
	}
	return num1 - num2, nil
}
func handleMultiplication(array []string) (float64, error) {
	num1, err := strconv.ParseFloat(array[0], 64)
	if err != nil {
		return 0.0, err
	}
	num2, err := strconv.ParseFloat(array[2], 64)
	if err != nil {
		return 0.0, err
	}
	return num1 * num2, nil
}
func handleDivision(array []string) (float64, error) {
	num1, err := strconv.ParseFloat(array[0], 64)
	if err != nil {
		return 0.0, err
	}
	num2, err := strconv.ParseFloat(array[2], 64)
	if err != nil {
		return 0.0, err
	}
	if num1 == 0 || num2 == 0 {
		return 0.0, errors.New("value can not be 0 in /")
	}
	return num1 / num2, nil
}

func HandleWriteRequestData(msg string) ([]byte, error) {
	arr := strings.Split(msg, " ")
	num1, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		return nil, err
	}
	num2, err := strconv.ParseFloat(arr[2], 64)
	if err != nil {
		return nil, err
	}
	requestData := model.Request{
		Num1:   num1,
		Num2:   num2,
		Action: arr[1],
	}
	var buf bytes.Buffer
	err = gob.NewEncoder(&buf).Encode(requestData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func HandleReadRequestData(b []byte) (*model.Request, error) {
	var requestData model.Request
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(&requestData)
	if err != nil {
		return nil, err
	}
	return &requestData, nil
}

func HandleWriteResponseData(result float64) ([]byte, error) {
	responseData := model.Response{
		Result: result,
	}
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(responseData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func HandleReadResponseData(b []byte) (*model.Response, error) {
	var response model.Response
	buf := bytes.NewBuffer(b)
	err := gob.NewDecoder(buf).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
