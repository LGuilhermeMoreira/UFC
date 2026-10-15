import os
import yaml # tive que instalar a lib pyyaml. doc: https://pypi.org/project/PyYAML/
import logging

class LogSystem:
    def __init__(self, level: str, log_file: str, format: str) -> None:
        self.level = level
        self.log_file = log_file
        self.format = format

        self.create_file()
        self.configure_logging()

    def create_file(self):
        if not os.path.exists(self.log_file):
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

def read_config(yaml_file: str):
    with open(yaml_file, "r") as file:
        config = yaml.safe_load(file)
    return config


