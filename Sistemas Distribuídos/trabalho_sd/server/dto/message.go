package dto

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Message struct {
	MessageType     uint            `json:"messageType"`
	RequestID       uuid.UUID       `json:"requestID"`
	ObjectReference string          `json:"objectReference"`
	Method          string          `json:"method"`
	Arguments       json.RawMessage `json:"arguments"`
	Error           any             `json:"error,omitempty"`
}

func (m Message) Debug() {
	if m.ObjectReference == "ping" {
		return
	}
	fmt.Printf("\n")
	fmt.Printf("MessageType: %v\n", m.MessageType)
	fmt.Printf("RequestID: %v\n", m.RequestID)
	fmt.Printf("ObjectReference: %v\n", m.ObjectReference)
	fmt.Printf("Method: %v\n", m.Method)
	var input map[string]any
	err := json.Unmarshal(m.Arguments, &input)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal dos argumentos", err)
	} else {
		fmt.Printf("Arguments: %v\n", input)
	}
	if m.Error != nil {
		fmt.Printf("Error: %v\n", m.Error)
	}
	fmt.Printf("\n")

}
