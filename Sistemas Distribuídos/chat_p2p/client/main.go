package main

import (
	handler "chat_p2p/client/handler"
	"net"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		handler.HandleRead(conn)
	}()

	go func() {
		defer wg.Done()
		handler.HandleWrite(conn)
	}()

	wg.Wait()
}
