package service

import (
	"trabalho_sd/contract"
	"trabalho_sd/dto"
	"trabalho_sd/models"
)

type Escola struct {
	BancoDeDados contract.BancoDeDadosInterface
}

func NewEscola(bancoDeDados contract.BancoDeDadosInterface) *Escola {
	return &Escola{BancoDeDados: bancoDeDados}
}

func (e *Escola) CadastrarAluno(aluno models.Aluno) (*dto.OutPutDefault, error) {
	err := e.BancoDeDados.CadastrarAluno(aluno)
	if err != nil {
		return nil, err
	}
	return &dto.OutPutDefault{
		Status:  201,
		Message: "Aluno cadastrado",
	}, nil
}

func (e *Escola) CadastrarProfessor(professor models.Professor) (*dto.OutPutDefault, error) {
	err := e.BancoDeDados.CadastrarProfessor(professor)
	if err != nil {
		return nil, err
	}
	return &dto.OutPutDefault{
		Status:  201,
		Message: "Professor cadastrado",
	}, nil

}

func (e *Escola) CadastrarDisciplina(disciplina models.Disciplina) (*dto.OutPutDefault, error) {
	err := e.BancoDeDados.CadastrarDisciplina(disciplina)
	if err != nil {
		return nil, err
	}
	return &dto.OutPutDefault{
		Status:  201,
		Message: "Disciplina cadastrado",
	}, nil
}

func (e *Escola) ListarTodasDisciplinas() (*dto.OutPutTodasDisciplinas, error) {
	data, err := e.BancoDeDados.ListarTodasDisciplinas()
	if err != nil {
		return nil, err
	}
	return &dto.OutPutTodasDisciplinas{
		Disciplinas: data,
		Status:      200,
	}, nil
}

func (e *Escola) BuscarAlunoPorCodigo(codigo string) (*dto.OutPutBuscaAlunoPorCodigo, error) {
	alunos, err := e.BancoDeDados.BuscarAlunoPorCodigo(codigo)
	if err != nil {
		return nil, err
	}

	return &dto.OutPutBuscaAlunoPorCodigo{
		Status: 200,
		Alunos: alunos,
	}, nil
}
