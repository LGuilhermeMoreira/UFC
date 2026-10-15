from pydantic import BaseModel

class CreateSapatoDto(BaseModel):
    modelo: str
    tamanho: int
    cor: str
    marca: str

class UpdateSapatoDto(BaseModel):
    modelo: str
    tamanho: int
    cor: str
    marca: str