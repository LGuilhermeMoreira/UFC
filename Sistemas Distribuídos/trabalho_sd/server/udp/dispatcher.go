package udp

import (
	"encoding/json"
	"log"
	"time"
	"trabalho_sd/dto"

	"github.com/google/uuid"
)

type Dispatcher struct {
	Skeleton  Skeleton
	mensagens map[uuid.UUID]dto.Message
}

func NewDispatcher(s Skeleton) *Dispatcher {
	return &Dispatcher{
		Skeleton:  s,
		mensagens: make(map[uuid.UUID]dto.Message),
	}
}

func (d *Dispatcher) Solve(buffer []byte) []byte {
	var msg dto.Message
	err := json.Unmarshal(buffer, &msg)
	if err != nil {
		log.Println("Erro ao deserializar a mensagem:", err)
		return []byte{}
	}
	msg.Debug()

	if resposta, loaded := d.mensagens[msg.RequestID]; loaded {
		log.Println("Mensagem duplicada detectada:", msg.RequestID)
		data, err := json.Marshal(resposta)
		if err != nil {
			log.Println("Erro ao serializar resposta do histórico:", err)
		}
		return data
	}

	var arguments any
	var ok bool

	switch msg.ObjectReference {
	case "Escola":
		arguments, ok = d.Skeleton.HandleRequest(&msg)
	case "ping":
		response := map[string]string{"msg": "pong"}
		data, _ := json.Marshal(response)
		return data
	default:
		arguments = map[string]any{
			"status": 404,
			"error":  "objectReference não encontrado",
		}
		ok = false
	}
	escreveMensagem(arguments, ok, &msg)
	msg.Debug()

	d.mensagens[msg.RequestID] = msg

	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("Erro ao serializar a mensagem:", err)
		return []byte{}
	}
	go func() {
		time.Sleep(20 * time.Second)
		delete(d.mensagens, msg.RequestID)
		log.Println("Mensagem removida do histórico:", msg.RequestID)
	}()
	return data
}

func escreveMensagem(data any, ok bool, message *dto.Message) {
	if !ok {
		message.Error = data
	} else {
		responseData, _ := json.Marshal(data)
		message.Arguments = responseData
	}
	message.MessageType = 1
}
