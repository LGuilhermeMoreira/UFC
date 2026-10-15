package handler

import (
	"encoding/json"
	"net"
	"sockets/entity"
	"time"
)

func Connection(conn net.Conn) {
	var req entity.Request
	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		json.NewEncoder(conn).Encode(entity.Response{
			Result: err.Error(),
		})
		return
	}
	var resp entity.Response
	switch req.Operation {
	case "**":
		resp = entity.GetResponse(entity.GetScienceCalculator(), req)
	case "@":
		resp = entity.GetResponse(entity.GetAnotherCalculator(), req)
	default:
		resp = entity.GetResponse(entity.GetSimpleCalculator(), req)
	}
	time.Sleep(100 * time.Millisecond)
	json.NewEncoder(conn).Encode(resp)
}
