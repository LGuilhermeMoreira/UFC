from model.sapato import Sapato
from uuid import UUID
import csv
import os
import hashlib #doc: https://docs.python.org/3/library/hashlib.html
from dto.dto import UpdateSapatoDto
import zipfile
from typing import Optional
from fastapi import HTTPException

class HandleFile():
    path : str

    def __init__(self,path : str) -> None:
        self.path = path
        self.createFile()

    def createFile(self) -> None:
        # configurando os campos no csv
        if not os.path.exists(self.path) or os.path.getsize(self.path) == 0:
            dados = ["id", "modelo", "tamanho", "cor", "marca", "created_at"]
            with open(self.path, mode='w', newline='', encoding='utf-8') as arquivo:
                fileWriter = csv.writer(arquivo)
                fileWriter.writerow(dados)
    
    def addSapato(self,sapato : Sapato) -> Sapato:
        dados = [sapato.id,sapato.modelo,sapato.tamanho,sapato.cor,sapato.marca,sapato.created_at]
        with open(self.path,mode='a',newline='',encoding='utf-8') as arquivo:
            fileWriter = csv.writer(arquivo)
            fileWriter.writerow(dados)
        return sapato

    def getSapatos(self) -> list:
        dados = []
        with open(self.path, mode='r', encoding='utf-8') as arquivo:
            leitor = csv.DictReader(arquivo)
            for linha in leitor:
                dados.append({
                    'id': linha['id'],
                    'modelo': linha['modelo'],
                    'tamanho': linha['tamanho'],
                    'cor': linha['cor'],
                    'marca': linha['marca'],
                    'created_at': linha['created_at'],
                })
        return dados
        

    def getSapatoById(self, id : UUID) -> dict:
        with open(self.path, mode='r', encoding='utf-8') as arquivo:
            leitor = csv.DictReader(arquivo)
            for linha in leitor:
               if linha['id'] == str(id):
                   return linha
        raise HTTPException(status_code=404,detail='ID não encontrado')

    def updateSapato(self,id : UUID,sapato : UpdateSapatoDto) -> dict:
        dados = []
        id_encontrado = False
        sapato_atualizado = None
        with open(self.path, mode='r', encoding='utf-8') as arquivo:
            leitor = csv.DictReader(arquivo)
            for linha in leitor:
                if linha['id'] == str(id):
                    linha['modelo'] = sapato.modelo
                    linha['tamanho'] = sapato.tamanho
                    linha['cor'] = sapato.cor
                    linha['marca'] = sapato.marca
                    id_encontrado = True
                    sapato_atualizado = linha
                dados.append(linha)
        
        if not id_encontrado:
            raise HTTPException(status_code=404,detail='ID não encontrado')
        with open(self.path, mode='w', newline='', encoding='utf-8') as arquivo:
            fileWriter = csv.DictWriter(arquivo, fieldnames=['id', 'modelo', 'tamanho', 'cor', 'marca', 'created_at'])
            fileWriter.writeheader()
            fileWriter.writerows(dados)
        return sapato_atualizado

    def deleteSapato(self,id : UUID) -> None:
        id_encontrado = False
        dados = []
        with open(self.path, mode='r', encoding='utf-8') as arquivo:
            leitor = csv.DictReader(arquivo)
            for linha in leitor:
                if linha['id'] != str(id):
                    dados.append(linha)
                else:
                    id_encontrado = True
        if not id_encontrado:
            raise HTTPException(status_code=404,detail='ID não encontrado')
        with open(self.path, mode='w', newline='', encoding='utf-8') as arquivo:
            fileWriter = csv.DictWriter(arquivo, fieldnames=['id', 'modelo', 'tamanho', 'cor', 'marca', 'created_at'])
            fileWriter.writeheader()
            fileWriter.writerows(dados)

    
    def getLenAllEntity(self) -> int:
        dados = []
        with open(self.path, mode='r', encoding='utf-8') as arquivo:
            leitor = csv.DictReader(arquivo)
            for linha in leitor:
                dados.append({
                    'id': linha['id'],
                    'modelo': linha['modelo'],
                    'tamanho': linha['tamanho'],
                    'cor': linha['cor'],
                    'marca': linha['marca'],
                    'created_at': linha['created_at'],
                })
        return {"len" : len(dados)}
    
    def getHashFile(self) -> str:
        m = hashlib.sha256()
        with open(self.path,mode='rb') as file: 
             for chunk in iter(lambda: file.read(4096), b""): 
                m.update(chunk)
        return {"hash" : m.hexdigest()}


    def createZip(self, zip_path: str) -> str:

        if not os.path.exists(self.path):
            raise FileNotFoundError("Arquivo CSV não encontrado.")

        with zipfile.ZipFile(zip_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
            zipf.write(self.path, arcname=os.path.basename(self.path))
        
        return zip_path

    def filterSapatos(
        self,
        modelo: Optional[str] = None,
        tamanho: Optional[int] = None,
        cor: Optional[str] = None,
        marca: Optional[str] = None,
    ) -> list:
        
        sapatos = self.getSapatos()
        
        if modelo:
            sapatos = [s for s in sapatos if s['modelo'].lower() == modelo.lower()]
        if tamanho:
            sapatos = [s for s in sapatos if int(s['tamanho']) == tamanho]
        if cor:
            sapatos = [s for s in sapatos if s['cor'].lower() == cor.lower()]
        if marca:
            sapatos = [s for s in sapatos if s['marca'].lower() == marca.lower()]

        return sapatos

