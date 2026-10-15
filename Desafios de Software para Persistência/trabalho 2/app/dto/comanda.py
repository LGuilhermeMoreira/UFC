from pydantic import BaseModel
from datetime import datetime

class ComandaCreate(BaseModel):
    id_cliente: int 
    id_mesa: int
    data_hora_abertura: datetime
    data_hora_fechamento: datetime | None
    status: str

class ComandaUpdate(BaseModel):
    id_cliente: int | None
    id_mesa: int | None
    data_hora_abertura: datetime | None
    data_hora_fechamento: datetime | None
    status: str | None