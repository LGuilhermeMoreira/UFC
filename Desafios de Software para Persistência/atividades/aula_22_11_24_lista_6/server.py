import socket
import pessoa_pb2

def create_person(name : str, age : int, email : str) -> str:
    person = pessoa_pb2.Person()
    person.name = name
    person.age = age 
    person.email = email
    return person.SerializeToString()
server = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
server.bind(('localhost', 5000))
server.listen(1)
print("Servidor aguardando conexão...")

conn, addr = server.accept()
print(f"Conexão de {addr}")

data = create_person('guilherme',20,'guilherme@email.com')
conn.sendall(data)
conn.close()
