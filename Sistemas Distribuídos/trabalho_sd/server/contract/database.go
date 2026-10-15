package contract

import "trabalho_sd/models"

type BancoDeDadosInterface interface {
	CadastrarAluno(models.Aluno) error
	CadastrarProfessor(models.Professor) error
	CadastrarDisciplina(models.Disciplina) error
	BuscarAlunoPorCodigo(string) ([]models.Aluno, error)
	ListarTodasDisciplinas() ([]models.Disciplina, error)
}
