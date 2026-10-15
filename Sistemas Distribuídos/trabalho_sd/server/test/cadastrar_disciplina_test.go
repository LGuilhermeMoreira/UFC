package test

import (
	"encoding/json"
	"net"
	"testing"
	"time"
	"trabalho_sd/dto"
	"trabalho_sd/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCadastrarDisciplina(t *testing.T) {
	cnfg := NewTestConnction()

	// arguemento que será passado para o servidor
	disciplina := dto.Disciplina{
		Nome:        "DSPP",
		Professores: []models.Professor{*models.NewProfessor("victor", "t@email,com")},
		Codigo:      "DSPP001",
	}

	data, err := json.Marshal(disciplina)
	if err != nil {
		t.Fatalf("Erro ao serializar 'codigo': %v", err)
	}

	// montando a mensagem.
	message := dto.Message{
		RequestID:       uuid.New(),
		MessageType:     0,
		ObjectReference: "Escola",
		Arguments:       data,
		Method:          "CadastrarDisciplina",
	}

	dataR, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Erro ao serializar 'message': %v", err)
	}

	_, err = cnfg.Conn.Write(dataR)
	if err != nil {
		t.Fatalf("Erro ao enviar mensagem: %v", err)
	}

	buffer := make([]byte, 1024)
	cnfg.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, _, err := cnfg.Conn.ReadFromUDP(buffer)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			t.Fatalf("Timeout ao receber resposta: %v", err)
		} else {
			t.Fatalf("Erro ao receber resposta: %v", err)
		}
	}

	var response dto.Message
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		t.Fatalf("Erro ao deserializar resposta: %v", err)
	}
	assert.NotEmpty(t, response)
}

func TestCadastrarDisciplinaFail(t *testing.T) {
	cnfg := NewTestConnction()

	// arguemento que será passado para o servidor
	disciplina := dto.Disciplina{
		Nome:        "DSPP",
		Professores: []models.Professor{*models.NewProfessor("victor", "t@email,com")},
		Codigo:      "DSPP001",
	}

	data, err := json.Marshal(disciplina)
	if err != nil {
		t.Fatalf("Erro ao serializar 'codigo': %v", err)
	}

	// montando a mensagem.
	message := dto.Message{
		RequestID:       uuid.New(),
		MessageType:     0,
		ObjectReference: "Escola",
		Arguments:       data,
		Method:          "CadastraDisciplina",
	}

	dataR, err := json.Marshal(message)
	if err != nil {
		t.Fatalf("Erro ao serializar 'message': %v", err)
	}

	_, err = cnfg.Conn.Write(dataR)
	if err != nil {
		t.Fatalf("Erro ao enviar mensagem: %v", err)
	}

	buffer := make([]byte, 1024)
	cnfg.Conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	n, _, err := cnfg.Conn.ReadFromUDP(buffer)

	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			t.Fatalf("Timeout ao receber resposta: %v", err)
		} else {
			t.Fatalf("Erro ao receber resposta: %v", err)
		}
	}

	var response dto.Message
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		t.Fatalf("Erro ao deserializar resposta: %v", err)
	}
	assert.NotEmpty(t, response.Error)
}
