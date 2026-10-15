import socket
from app.models import *
from app.connection import UDPCliente

class EscolaService:
    def __init__(self,udpCliente:UDPCliente):
        self.udp = udpCliente
        

    def _processa_requisicao(self, service: str, method: str, arguments) -> Message | None:
        msg = Message(service, method, arguments.__dict__)
        msg_json = msg.to_json()
        tentativas = 3

        for _ in range(tentativas):
            try:
                self.udp.enviar_mensagem(msg_json)
                data = self.udp.receber_mensagem()
                return Message.from_json(data)
            except socket.timeout:
                print(f"Tentativa falhou, reenviando... ({_ + 1}/{tentativas})")
            except Exception as e:
                print(f"Erro ao processar requisição: {e}")
                return None
        
        return None

    def cadastrar_aluno(self, aluno: Aluno) -> Message | None:
        return self._processa_requisicao("Escola", "CadastrarAluno", aluno)

    def cadastrar_disciplina(self, disciplina: Disciplina) -> Message | None:
        return self._processa_requisicao("Escola", "CadastrarDisciplina", disciplina)

    def cadastrar_professor(self, professor: Professor) -> Message | None:
        return self._processa_requisicao("Escola", "CadastrarProfessor", professor)

    def buscar_aluno_por_codigo(self, codigo: Codigo) -> Message | None:
        return self._processa_requisicao("Escola", "BuscarAlunoPorCodigo", codigo)
    
    def listar_todas_disciplinas(self,default : Default):
        return self._processa_requisicao("Escola","ListarTodasDisciplinas",default)
        
