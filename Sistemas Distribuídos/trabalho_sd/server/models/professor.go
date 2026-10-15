package models

type Professor struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
}

func NewProfessor(nome, email string) *Professor {
	return &Professor{
		Nome:  nome,
		Email: email,
	}
}
