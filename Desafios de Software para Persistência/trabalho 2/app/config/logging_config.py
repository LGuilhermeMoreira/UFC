import logging
import os
from logging.handlers import RotatingFileHandler

def setup_logging():
    log_dir = "logs"  
    os.makedirs(log_dir, exist_ok=True)
    log_file = os.path.join(log_dir, "api.log")  
    
    
    log_format = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')

   
    file_handler = RotatingFileHandler(
        log_file,
        maxBytes=10 * 1024 * 1024,  
        backupCount=5  
    )
    file_handler.setFormatter(log_format)

    
    logger = logging.getLogger("api_logger")
    logger.setLevel(logging.DEBUG)  
    logger.addHandler(file_handler)

    return logger