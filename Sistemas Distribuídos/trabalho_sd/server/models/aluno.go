package models

import "github.com/google/uuid"

type Aluno struct {
	Nome        string       `json:"nome"`
	Matricula   string       `json:"matricula"`
	Curso       string       `json:"curso"`
	Cpf         string       `json:"cpf"`
	Disciplinas []Disciplina `json:"disciplinas"`
}

func NewAluno(nome, curso, cpf string, disciplinas ...Disciplina) *Aluno {
	matricula := uuid.New().String()
	return &Aluno{
		Nome:        nome,
		Matricula:   matricula,
		Curso:       curso,
		Cpf:         cpf,
		Disciplinas: disciplinas,
	}
}
