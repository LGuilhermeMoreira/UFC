from pydantic import BaseModel

class ClienteCreate(BaseModel):
    nome: str
    telefone: str | None  
    email:  str |  None
    cpf: str

class ClienteUpdate(BaseModel):
    nome: str | None
    telefone: str | None
    email: str | None
    cpf: str | None