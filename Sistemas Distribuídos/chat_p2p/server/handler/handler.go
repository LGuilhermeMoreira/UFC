package handler

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	go HandleRead(conn)
	HandleWrite(conn)
}

func HandleRead(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from socket:", err)
			return
		}
		if n > 0 {
			log.Println("Received:", string(buffer[:n]))
		}
	}
}

func HandleWrite(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading input from keyboard:", err)
			continue
		}
		msg = strings.TrimSpace(msg)
		if msg == "EXIT" {
			fmt.Println("Exiting server.")
			os.Exit(0)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Println("Error writing to socket:", err)
		}
	}
}
