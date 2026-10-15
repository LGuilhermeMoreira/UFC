package entity

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/exp/slices"
)

const port = ":8080"

type User struct {
	input  string
	Result Response
}

var once = sync.Once{}
var singletonUser *User

func GetUser() *User {
	once.Do(func() {
		singletonUser = &User{
			input: "10 * 10",
		}
	})
	return singletonUser
}

func (u *User) SendRequest() error {
	conn, err := net.Dial("tcp", port)
	if err != nil {
		return err
	}
	defer conn.Close()
	req, err := inputToRequest(u.input)
	if err != nil {
		return err
	}
	res, err := handleWriteRead(conn, *req)
	if err != nil {
		return err
	}
	u.Result = *res
	return nil
}

func (u *User) ShowResponse() {
	fmt.Println(u.Result.Result)
}

func handleWriteRead(conn net.Conn, req Request) (*Response, error) {
	err := json.NewEncoder(conn).Encode(req)
	if err != nil {
		return nil, err
	}
	var res Response
	err = json.NewDecoder(conn).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func validateInput(input string) bool {
	result := strings.Split(input, " ")
	if len(result) != 3 {
		return false
	}
	operation := result[1]
	symbols := []string{"+", "-", "*", "/", "**", "@"}
	contains := slices.Contains(symbols, operation)
	if !contains {
		return false
	}
	return true
}

func inputToRequest(input string) (*Request, error) {
	isValid := validateInput(input)
	if !isValid {
		return nil, errors.New("invalid input")
	}
	result := strings.Split(input, " ")
	num1, err := strconv.ParseFloat(result[0], 64)
	if err != nil {
		return nil, err
	}
	num2, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return nil, err
	}
	return &Request{
		Num1:      num1,
		Num2:      num2,
		Operation: result[1],
	}, nil
}
