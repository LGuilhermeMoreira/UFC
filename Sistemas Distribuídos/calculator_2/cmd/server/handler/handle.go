package handler

import (
	"encoding/json"
	"net"
	"sockets/entity"
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
		resp = entity.GetResponse(entity.ScienceCalculator{}, req)
	case "@":
		resp = entity.GetResponse(entity.AnotherCalculator{}, req)
	default:
		resp = entity.GetResponse(entity.SimpleCalculator{}, req)
	}
	json.NewEncoder(conn).Encode(resp)
	return
}
