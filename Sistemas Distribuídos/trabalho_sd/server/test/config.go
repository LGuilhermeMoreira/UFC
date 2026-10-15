package test

import "net"

type TestConnection struct {
	Address *net.UDPAddr
	Conn    *net.UDPConn
}

func NewTestConnction() *TestConnection {
	addr, err := net.ResolveUDPAddr("udp", "localhost:4567")
	if err != nil {
		return nil
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil
	}
	return &TestConnection{
		Address: addr,
		Conn:    conn,
	}
}
