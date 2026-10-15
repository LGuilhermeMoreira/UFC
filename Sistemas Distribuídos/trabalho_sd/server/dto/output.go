package dto

import "trabalho_sd/models"

type OutPutBuscaAlunoPorCodigo struct {
	Alunos []models.Aluno `json:"alunos"`
	Status int            `json:"status"`
}

type OutPutDefault struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type OutPutTodasDisciplinas struct {
	Status      int                 `json:"status" `
	Disciplinas []models.Disciplina `json:"disciplinas"`
}
