package main

import (
	"trabalho_sd/db"
	"trabalho_sd/service"
	"trabalho_sd/udp"
)

func main() {
	// DI
	db := db.NewBancoDeDados()
	escola := service.NewEscola(db)
	skeleton := udp.NewSkeleton(*escola)
	dispatcher := udp.NewDispatcher(*skeleton)
	serv := udp.UDPServer{
		Port:       ":4567",
		Ip:         "localhost",
		ServerType: "udp",
		Dispatcher: *dispatcher,
	}
	serv.Start()
}
