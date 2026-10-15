from dataclasses import dataclass

@dataclass
class Aluno:
    nome : str
    curso : str
    cpf : str
    disciplinas : any