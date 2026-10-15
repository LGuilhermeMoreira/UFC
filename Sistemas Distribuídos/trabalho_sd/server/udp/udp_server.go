package udp

import (
	"log"
	"math/rand"
	"net"
)

type UDPServer struct {
	Port       string
	Ip         string
	ServerType string
	Dispatcher Dispatcher
}

func (s *UDPServer) Start() {
	// criando
	addr, err := net.ResolveUDPAddr(s.ServerType, s.Ip+s.Port)
	if err != nil {
		log.Fatalf("Error creating addr: %s", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
	defer conn.Close()

	log.Printf("Server listening on port %v\n", s.Port)
	for {
		buffer := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Received message from %v\n", addr)
		data := s.Dispatcher.Solve(buffer[:n])
		randNum := rand.Int()
		// && randNum%3 == 0 && randNum%5 == 0
		if randNum%2 == 0 {
			log.Println("Simulando perda de pacote (não enviando a resposta):", randNum)
		} else {
			conn.WriteToUDP(data, addr)
		}
	}

}
