package dto

import (
	"trabalho_sd/models"
)

type Aluno struct {
	Nome        string              `json:"nome"`
	Curso       string              `json:"curso"`
	Cpf         string              `json:"cpf"`
	Disciplinas []models.Disciplina `json:"disciplinas"`
}
