package handler

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func HandleRead(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading message:", err)
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
			log.Println("Error reading from keyboard:", err)
			continue
		}
		msg = strings.TrimSpace(msg)
		if msg == "EXIT" {
			fmt.Println("Exiting chat.")
			os.Exit(0)
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Println("Error writing message:", err)
		}
	}
}
