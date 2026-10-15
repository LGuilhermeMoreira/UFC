#python -m uvicorn main:app --reload
from fastapi import FastAPI
from typing import Union

app = FastAPI()

@app.get("/")
async def hello_world():
    return {"Hello":"world"}

@app.get('/item/{item_id}')
async def get_item(item_id : int, nome : Union[str,None] = None):
    return {"item_id": item_id, "nome" : nome}

#python -m uvicorn main:app --reload
