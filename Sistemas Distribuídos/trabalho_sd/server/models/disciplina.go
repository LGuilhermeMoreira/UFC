package models

type Disciplina struct {
	Nome        string      `json:"nome"`
	Codigo      string      `json:"codigo"`
	Professores []Professor `json:"professores"`
}

func NewDisciplina(nome, codigo string, professores ...Professor) *Disciplina {
	return &Disciplina{
		Nome:        nome,
		Codigo:      codigo,
		Professores: professores,
	}
}
