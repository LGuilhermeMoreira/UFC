package main

import (
	"log"
	"net"
	"sockets/handler"
)

const (
	Type = "tcp"
	Port = ":8080"
)

func main() {
	srv, err := net.Listen(Type, Port)
	if err != nil {
		log.Fatal(err)
	}
	defer srv.Close()
	for {
		conn, err := srv.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handler.Connection(conn)
	}
}
