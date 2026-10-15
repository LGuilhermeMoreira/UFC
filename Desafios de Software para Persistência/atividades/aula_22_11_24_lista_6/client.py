import pessoa_pb2
import socket

client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client.connect(('localhost', 5000))

data = client.recv(1024)

person = pessoa_pb2.Person()
person.ParseFromString(data)

print(f"Recebido: {person.name}, {person.age}, {person.email}")  # Isso deve funcionar
client.close()
