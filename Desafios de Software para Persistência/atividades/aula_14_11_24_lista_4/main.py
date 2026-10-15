from livro import Book
from xml_file import XMLFileParse
from fastapi import FastAPI,HTTPException

xml = XMLFileParse('file.xml')
app = FastAPI()

@app.post('/book/', status_code=201)
async def CreateBook(book : Book):
    try:
        xml.Save(book)
        return book
    except Exception as e:
        raise HTTPException(500,detail=f'erro no servidor: {e}')

@app.get('/book/',status_code=200 )
async def GetAllBooks():
    try:
        response =  xml.FindAll()
        return response
    except Exception as e:
       raise HTTPException(500,detail=f'erro no servidor: {e}')

@app.get('/book/{id}',status_code=200 )
async def GetAllBooks(id):
    try:
        response =  xml.FindByID(id)
        return response
    except Exception as e:
       raise HTTPException(500,detail=f'erro no servidor: {e}')

@app.put('/book/{id}',status_code=200)
async def Update(id,book : Book):
    try:
        response =  xml.Update(id,book)
        return response
    except Exception as e:
       raise HTTPException(500,detail=f'erro no servidor: {e}')

@app.delete('/book/{id}',status_code=204)
async def Delete(id):
    try:
        response =  xml.Delete(id)
        return response
    except Exception as e:
       raise HTTPException(500,detail=f'erro no servidor: {e}')
