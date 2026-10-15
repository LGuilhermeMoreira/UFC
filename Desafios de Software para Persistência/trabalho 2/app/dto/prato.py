from pydantic import BaseModel
class PratoCreate(BaseModel):
    nome: str
    preco: float
    descricao:  str | None
    disponivel: bool | None

class PratoUpdate(BaseModel):
    nome:  str | None
    preco: float | None
    descricao:  str | None
    disponivel: bool | None
