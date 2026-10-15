package main

import (
	"net"
	"sockets/handler"
)

func main() {
	list, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := list.Accept()
		if err != nil {
			panic(err)
		}
		go handler.HandleCalculateWithStruct(conn)
	}
}
