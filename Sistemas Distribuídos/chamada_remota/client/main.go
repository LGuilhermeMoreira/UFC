package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Proxy struct
type Proxy struct {
	serverAddr string
}

func (p *Proxy) sendRequest(request string) string {
	conn, err := net.Dial("tcp", p.serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return ""
	}
	defer conn.Close()

	// Send request
	fmt.Fprintln(conn, request)

	// Receive response
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return ""
	}

	return response
}

func main() {
	proxy := &Proxy{serverAddr: "localhost:8080"}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Choose operation (add, sub, mul, div) or 'exit':")
		scanner.Scan()
		operation := scanner.Text()
		if operation == "exit" {
			break
		}

		fmt.Println("Enter two numbers (space-separated):")
		scanner.Scan()
		params := scanner.Text()

		request := fmt.Sprintf("%s %s", operation, params)
		response := proxy.sendRequest(request)
		fmt.Println("Response from server:", response)
	}
}
