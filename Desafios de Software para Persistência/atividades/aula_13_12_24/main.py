from database import get_db
from fastapi import FastAPI,Depends,HTTPExecption

app = FastAPI()

@app.post("/alunos")
def criar_aluno(nome:str,email:str,db:Session = Depends(get_db)):
    aluno = Aluno(nome=nome,email=email)
    db.add(aluno)
