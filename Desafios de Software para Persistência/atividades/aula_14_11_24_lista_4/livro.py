from pydantic import BaseModel

class Book(BaseModel):
    id : int | None = None # para se ajustar a rota de update
    titulo: str
    autor : str
    ano : int
    genero: str

    def validateToSave(self):
        if self.id is None:
            raise Exception('invalid book, id is missing')
        else:
            return True
    
