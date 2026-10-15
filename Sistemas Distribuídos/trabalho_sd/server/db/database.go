package db

import (
	"errors"
	"fmt"
	"trabalho_sd/models"
)

type BancoDeDados struct {
	Alunos      []models.Aluno
	Disciplinas []models.Disciplina
	Professores []models.Professor
}

func NewBancoDeDados() *BancoDeDados {
	return &BancoDeDados{
		Professores: []models.Professor{
			*models.NewProfessor("samy", "samy@mail.com"),
			*models.NewProfessor("marcos", "marcos@email.com"),
			*models.NewProfessor("thigas", "thigas@email.com"),
			*models.NewProfessor("livia", "livia@email.com"),
		},
		Disciplinas: []models.Disciplina{
			*models.NewDisciplina("SD", "SD001", *models.NewProfessor("marcos", "marcos@email.com")),
			*models.NewDisciplina("SO", "SO001", *models.NewProfessor("thigas", "thigas@email.com")),
			*models.NewDisciplina("BD", "BD001", *models.NewProfessor("livia", "livia@email.com")),
			*models.NewDisciplina("IA", "IA001", *models.NewProfessor("samy", "samy@mail.com")),
		},
		Alunos: []models.Aluno{
			*models.NewAluno("Guigui", "CC", "1233212221",
				*models.NewDisciplina("SD", "SD001",
					*models.NewProfessor("marcos", "marcos@email.com"))),

			*models.NewAluno("Iaia Pirata", "CC", "1233212225",
				*models.NewDisciplina("SD", "SD001",
					*models.NewProfessor("marcos", "marcos@email.com"))),

			*models.NewAluno("Kaynara", "CC", "1233212222",
				*models.NewDisciplina("SO", "SO001",
					*models.NewProfessor("thigas", "thigas@email.com"))),

			*models.NewAluno("Cabeça", "CC", "1233212223",
				*models.NewDisciplina("BD", "BD001",
					*models.NewProfessor("livia", "livia@email.com"))),

			*models.NewAluno("jss de CC", "CC", "1233212224",
				*models.NewDisciplina("IA", "IA001",
					*models.NewProfessor("samy", "samy@mail.com"))),
		},
	}
}

func (b *BancoDeDados) CadastrarAluno(aluno models.Aluno) error {
	fmt.Printf("Nome:%v\tMatricula:%v\tCurso:%v\tCpf:%v\n", aluno.Nome, aluno.Matricula, aluno.Curso, aluno.Cpf)
	for _, alunoList := range b.Alunos {
		if alunoList.Cpf == aluno.Cpf {
			return errors.New("cpf já atribuido a outro aluno")
		}
	}
	b.Alunos = append(b.Alunos, aluno)
	return nil
}

func (b *BancoDeDados) CadastrarDisciplina(disciplina models.Disciplina) error {
	fmt.Printf("Nome: %v\tCodigo %v\n", disciplina.Nome, disciplina.Codigo)
	for _, disciplinaList := range b.Disciplinas {
		if disciplinaList.Codigo == disciplina.Codigo {
			return fmt.Errorf("código já está atribuido a disciplina: %v", disciplinaList.Nome)
		}
	}
	b.Disciplinas = append(b.Disciplinas, disciplina)
	return nil
}

func (b *BancoDeDados) CadastrarProfessor(professor models.Professor) error {
	fmt.Printf("Nome:%v\tEmail:%v\n", professor.Nome, professor.Email)
	for _, professorList := range b.Professores {
		if professorList.Email == professor.Email {
			return fmt.Errorf("email já pertencente ao professor: %s", professorList.Nome)
		}
	}
	b.Professores = append(b.Professores, professor)
	return nil
}

func (b *BancoDeDados) BuscarAlunoPorCodigo(codigo string) ([]models.Aluno, error) {
	var responseData []models.Aluno
	for _, aluno := range b.Alunos {
		for _, disciplina := range aluno.Disciplinas {
			if disciplina.Codigo == codigo {
				responseData = append(responseData, aluno)
			}
		}
	}
	if len(responseData) == 0 {
		return nil, errors.New("nenhum aluno encontrado")
	}
	return responseData, nil
}

func (b *BancoDeDados) ListarTodasDisciplinas() ([]models.Disciplina, error) {
	if len(b.Disciplinas) == 0 {
		return nil, errors.New("nenhuma disciplina cadastrada")
	}
	return b.Disciplinas, nil
}
