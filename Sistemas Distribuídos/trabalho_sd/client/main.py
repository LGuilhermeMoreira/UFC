from app.utils import Interface
from app.proxy import EscolaService
from app.connection import UDPCliente

#DI
try:    
    udp = UDPCliente(debug=False)
    es = EscolaService(udpCliente=udp)
    interface = Interface(escolaService=es)
    interface.start()
except Exception as e:
    print(f"Erro no programa: {e}")

#start