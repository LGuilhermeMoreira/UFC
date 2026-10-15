from pydantic import BaseModel
class IngredienteCreate(BaseModel):
    nome: str
    estoque: bool
    quantidade_estoque: float
    peso: float
    cor : str

class IngredienteUpdate(BaseModel):
    nome: str | None
    estoque:bool | None
    quantidade_estoque: float | None
    peso: float | None
    cor : str | None
