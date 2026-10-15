import os
import xml.etree.ElementTree as ET
from livro import Book

class XMLFileParse:
    def __init__(self, pathToFile: str):
        self.pathToFile = pathToFile
        
        if not os.path.exists(self.pathToFile):
            self.create_empty_xml()

    def create_empty_xml(self):
        root = ET.Element('library') 
        tree = ET.ElementTree(root)
        tree.write(self.pathToFile, encoding='utf-8', xml_declaration=True)

    def load_xml(self):
        return ET.parse(self.pathToFile)

    def FindAll(self):
        tree = self.load_xml()
        root = tree.getroot()
        books = []
        for book_elem in root.findall('book'):
            books.append(Book(
                id=int(book_elem.get('id')),
                titulo=book_elem.find('title').text,
                autor=book_elem.find('author').text,
                ano=int(book_elem.find('year').text),
                genero=book_elem.find('genre').text
            ))
        return books or []  # Retorna uma lista vazia se não houver livros

    def Delete(self, id: int):
        tree = self.load_xml()
        root = tree.getroot()
        for book_elem in root.findall('book'):
            if book_elem.get('id') == str(id):
                root.remove(book_elem)
                tree.write(self.pathToFile, encoding='utf-8', xml_declaration=True)
                return True
        return False  # Indica que o ID não foi encontrado


    def FindByID(self, id: int):
        tree = self.load_xml()
        root = tree.getroot()
        for book_elem in root.findall('book'):
            if book_elem.get('id') == str(id):
                return Book(
                    id=int(book_elem.get('id')),
                    titulo=book_elem.find('title').text,
                    autor=book_elem.find('author').text,
                    ano=int(book_elem.find('year').text),
                    genero=book_elem.find('genre').text
                )
        return None  

    def Save(self, book: Book):
        tree = self.load_xml()
        root = tree.getroot()
        book.validateToSave()

        book_elem = ET.SubElement(root, 'book', id=str(book.id))
        ET.SubElement(book_elem, 'title').text = book.titulo
        ET.SubElement(book_elem, 'author').text = book.autor
        ET.SubElement(book_elem, 'year').text = str(book.ano)
        ET.SubElement(book_elem, 'genre').text = book.genero

        tree.write(self.pathToFile, encoding='utf-8', xml_declaration=True)

    def Update(self, id: int, book: Book):
        tree = self.load_xml()
        root = tree.getroot()
        for book_elem in root.findall('book'):
            if book_elem.get('id') == str(id):
                book_elem.find('title').text = book.titulo
                book_elem.find('author').text = book.autor
                book_elem.find('year').text = str(book.ano)
                book_elem.find('genre').text = book.genero
                book.id = int(id)
                tree.write(self.pathToFile, encoding='utf-8', xml_declaration=True)
                return book
        return {}

