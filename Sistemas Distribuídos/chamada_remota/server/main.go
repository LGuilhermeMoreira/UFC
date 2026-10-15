package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Calc struct implements the actual operations
type Calc struct{}

func (c *Calc) Add(a, b float64) float64 {
	return a + b
}

func (c *Calc) Sub(a, b float64) float64 {
	return a - b
}

func (c *Calc) Mul(a, b float64) float64 {
	return a * b
}

func (c *Calc) Div(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

// Skeleton struct
type Skeleton struct {
	calc *Calc
}

func (s *Skeleton) handleRequest(request string) string {
	parts := strings.Fields(request)
	if len(parts) < 3 {
		return "Invalid request"
	}

	op := parts[0]
	a, err1 := strconv.ParseFloat(parts[1], 64)
	b, err2 := strconv.ParseFloat(parts[2], 64)

	if err1 != nil || err2 != nil {
		return "Invalid numbers"
	}

	var result float64
	switch op {
	case "add":
		result = s.calc.Add(a, b)
	case "sub":
		result = s.calc.Sub(a, b)
	case "mul":
		result = s.calc.Mul(a, b)
	case "div":
		result = s.calc.Div(a, b)
	default:
		return "Unknown operation"
	}

	return fmt.Sprintf("%.2f", result)
}

// Dispatcher function
func startServer(address string) {
	calc := &Calc{}
	skeleton := &Skeleton{calc: calc}

	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server is running on", address)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, skeleton)
	}
}

func handleConnection(conn net.Conn, skeleton *Skeleton) {
	defer conn.Close()

	request, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	response := skeleton.handleRequest(strings.TrimSpace(request))
	fmt.Fprintln(conn, response)
}

func main() {
	startServer("localhost:8080")
}
