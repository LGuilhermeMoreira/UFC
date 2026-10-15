import os
import json
import yaml # tive que instalar a lib pyyaml. doc: https://pypi.org/project/PyYAML/
import logging

class LogSystem:
    def __init__(self, level: str, log_file: str, format: str, json_file: str) -> None:
        self.level = level
        self.log_file = log_file
        self.format = format
        self.json_file = json_file

        self.create_file()
        self.configure_logging()

    def create_file(self):
        with open(self.log_file, mode='w') as arquivo:
            pass

    def configure_logging(self):
        logging.basicConfig(
            level=getattr(logging, self.level.upper(), logging.INFO),
            filename=self.log_file,
            format=self.format
        )

    def add_log(self, message: str, level: str = "INFO"):
        logger = logging.getLogger()
        log_function = getattr(logger, level.lower(), logger.info)
        log_function(message)

    def process_json(self):
        if not os.path.exists(self.json_file):
            self.add_log(f"Arquivo JSON '{self.json_file}' não encontrado.", "ERROR")
            return

        try:
            with open(self.json_file, "r") as arquivo:
                data = json.load(arquivo)
                self.add_log(f"Arquivo JSON '{self.json_file}' carregado com sucesso.", "INFO")
        except json.JSONDecodeError as e:
            self.add_log(f"Erro ao carregar o arquivo JSON: {e}", "ERROR")
            return

        for record in data:
            if record.get("id") and record.get("name") and isinstance(record.get("age"), int):
                self.add_log(f"Processando registro: {record}", "INFO")
            else:
                self.add_log(f"Dado inválido: {record}", "WARNING")


def read_config(yaml_file: str):
    with open(yaml_file, "r") as file:
        config = yaml.safe_load(file)
    return config



config = read_config("config.yml")

log_config = config["logging"]
data_config = config["data"]

log_system = LogSystem(
    level=log_config["level"],
    log_file=log_config["file"],
    format=log_config["format"],
    json_file=data_config["file"]
)

log_system.process_json()
