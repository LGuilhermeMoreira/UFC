from pydantic import BaseModel
class MesaCreate(BaseModel):
    numero_mesa: int
    capacidade: int
    ocupada: bool
    numero_pessoas: int | None

class MesaUpdate(BaseModel):
    numero_mesa: int | None
    capacidade: int | None
    ocupada: bool | None
    numero_pessoas: int | None