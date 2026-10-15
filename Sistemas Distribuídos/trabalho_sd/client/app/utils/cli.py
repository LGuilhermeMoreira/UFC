from app.models import *
from app.proxy import *
import os
class Interface:
    def __init__(self,escolaService:EscolaService):
        self.es = escolaService
    
    def start(self):
        es = self.es
        while True:
            print("\nEscolha uma ação:")
            print("1. Buscar Aluno por Código")
            print("2. Cadastrar Aluno")
            print("3. Cadastrar Disciplina")
            print("4. Cadastrar Professor")
            print("5. Listar Disciplinas")
            print("0. Sair")

            choice = input("Opção: ")

            if choice == '1':
                codigo_str = input("Digite o código da disciplina: ")
                codigo = Codigo(codigo=codigo_str)
                self.limpar_terminal()
                msg = es.buscar_aluno_por_codigo(codigo)
                if msg is not None:
                    self.handle_response(msg)
            
            elif choice == '2':
                nome = input("Nome do aluno: ")
                curso = input("Curso do aluno: ")
                cpf = input("CPF do aluno: ")
                prof = Professor("marcos","marcos@email.com").__dict__
                disciplinas = [Disciplina("SD","SD001",prof).__dict__]
                self.limpar_terminal()
                aluno = Aluno(nome=nome, curso=curso, cpf=cpf, disciplinas=disciplinas)
                msg = es.cadastrar_aluno(aluno=aluno)
                if msg is not None:
                    self.handle_response(msg)
               
            elif choice == '3':
                nome = input("Nome da disciplina: ")
                codigo = input("Código da disciplina: ")
                professor = Professor("vitu","vitu@mail.com").__dict__
                professores = []
                professores.append(professor)
                self.limpar_terminal()
                disciplina = Disciplina(nome, codigo,professores)
                msg = es.cadastrar_disciplina(disciplina=disciplina)
                if msg is not None:
                    self.handle_response(msg)
            
            elif choice == '4':
                nome = input("Nome do professor: ")
                email = input("Email do professor: ")
                professor = Professor(nome=nome, email=email)
                self.limpar_terminal()
                msg = es.cadastrar_professor(professor=professor)
                if msg is not None:
                    self.handle_response(msg)
                    
            elif choice == '5':
                self.limpar_terminal()
                dft = Default("Quero todas as disciplinas")
                msg = es.listar_todas_disciplinas(dft)
                if msg is not None:
                    self.handle_response(msg)
            elif choice == '0':
                break
            
            else:
                print("Opção inválida. Tente novamente.")
                
    
    def handle_response(self,msg : Message):
        if msg.error is not None:
            status = msg.error.get("status")
            error = msg.error.get("error")
            print(f"Status: {status}")
            print(f"error: {error}")
            return
        match msg.method:
            case "CadastrarAluno" | "CadastrarDisciplina" | "CadastrarProfessor":
                status = msg.arguments.get("status")
                mensagem = msg.arguments.get("message")
                print(f"Status: {status}")
                print(f"mensagem: {mensagem}")
            case "BuscarAlunoPorCodigo":
                status = msg.arguments.get("status")
                alunos = msg.arguments.get("alunos", [])
                print(f"Status: {status}")
                for aluno in alunos:
                    print(f"Nome: {aluno['nome']}, CPF: {aluno['cpf']}")
            case "ListarTodasDisciplinas":
                status = msg.arguments.get("status")  
                disciplinas = msg.arguments.get("disciplinas", [])  
                print(f"Status: {status}")
                for disciplina in disciplinas:
                    nome = disciplina.get("nome")  
                    codigo = disciplina.get("codigo")
                    professores = disciplina.get("professores", []) 
                    
                    for professor in professores:
                        nome_professor = professor.get("nome")  
                        email_professor = professor.get("email")  
                        print(f"Disciplina: {nome}, Código: {codigo}, Professor: {nome_professor}, Email: {email_professor}")
            case _:
                print("Erro")
            
                
    def limpar_terminal(self):
        os.system('cls' if os.name == 'nt' else 'clear')
