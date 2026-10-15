import socket
from app.models import *
class UDPCliente:
    
    def __init__(self,debug : bool):
        self.addr = ('localhost', 4567)
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self.socket.settimeout(1) 
        
        ok = self._ping_no_servidor(debug)
        if not ok:
            raise ConnectionError("Não foi possível se conectar ao servidor")
        
    def enviar_mensagem(self,message_json: str = ""):    
        self.socket.sendto(message_json.encode('utf-8'),self.addr)
        
    def receber_mensagem(self):
        data, _ = self.socket.recvfrom(1024)
        return data

    def _ping_no_servidor(self,debug):
            if debug:
               print("enviando: ")
               print({"msg" : "ping"}) 
            msg = Message(objectReference="ping",method="ping",arguments={"msg" : "ping"})
            data = msg.to_json()
            self.enviar_mensagem(data)
            data = self.receber_mensagem()
            if debug:
                print("recebendo: ")
                print(str(data))
            if data is None:
                return False
            return True