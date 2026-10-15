package dto

import "trabalho_sd/models"

type Disciplina struct {
	Nome        string             `json:"nome"`
	Codigo      string             `json:"codigo"`
	Professores []models.Professor `json:"professores"`
}
