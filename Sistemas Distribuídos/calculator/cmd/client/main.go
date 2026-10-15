package main

import (
	"flag"
	"fmt"
	"net"
	"sockets/handler"
)

func main() {
	port := flag.String("port", "8080", "The port to connect to")
	msg := flag.String("msg", "hello", "The message to send")
	flag.Parse()
	conn, err := net.Dial("tcp", fmt.Sprintf(":%s", *port))
	fmt.Println(*msg)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	request, err := handler.HandleWriteRequestData(*msg)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(request)
	if err != nil {
		panic(err)
	}
	responseData := make([]byte, 1024)
	n, err := conn.Read(responseData)
	if err != nil {
		panic(err)
	}
	response, err := handler.HandleReadResponseData(responseData[:n])
	if err != nil {
		panic(err)
	}
	response.Show()
}
