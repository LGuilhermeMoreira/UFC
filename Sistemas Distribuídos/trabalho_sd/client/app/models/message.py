import json
import uuid
from dataclasses import dataclass, field

@dataclass
class Message:
    objectReference: str
    method: str
    arguments: any
    messageType: int = field(default=0)
    requestID: str = field(default_factory=lambda: str(uuid.uuid4()))
    error: any = None
 
    def to_json(self):
        return json.dumps({
            "messageType": self.messageType,
            "requestID": self.requestID,
            "objectReference": self.objectReference,
            "method": self.method,
            "arguments": self.arguments,  
            "error": self.error
        })

    @classmethod
    def from_json(cls, json_str):
        try:
            data = json.loads(json_str)
            return cls(**data)
        except json.JSONDecodeError:
            print("Erro: JSON inválido")
            return None

