package test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
	"trabalho_sd/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// 90784ddf-e491-4433-84ab-eb0d973b4c6e uuid
func TestMensagensDuplicadas(t *testing.T) {
	err1 := send()
	assert.Nil(t, err1)
	err2 := send()
	assert.Nil(t, err2)
}

func send() error {
	cnfg := NewTestConnction()

	aluno := dto.Aluno{
		Nome:  "A",
		Curso: "B",
		Cpf:   "12332100099",
	}

	data, err := json.Marshal(aluno)
	if err != nil {
		return err
	}

	// montando a mensagem.
	id, err := uuid.Parse("90784ddf-e491-4433-84ab-eb0d973b4c6e")
	if err != nil {
		return err
	}
	message := dto.Message{
		RequestID:       id,
		MessageType:     0,
		ObjectReference: "Escola",
		Arguments:       data,
		Method:          "CadastrarAluno",
	}

	dataR, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = cnfg.Conn.Write(dataR)
	if err != nil {
		return err
	}

	buffer := make([]byte, 1024)
	cnfg.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, _, err := cnfg.Conn.ReadFromUDP(buffer)

	if err != nil {
		return err
	}

	var response dto.Message
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		return err
	}

	if response.Error == nil {
		return nil
	} else {
		return errors.New("erro no campo Error de response")
	}
}
