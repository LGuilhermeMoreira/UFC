from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlmodel import SQLModel
from .config import settings
import sys
try:
    engine = create_engine(settings.DATABASE_URL)
    SQLModel.metadata.create_all(engine)
    SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)
except:
    print("banco de dados não conectado")
    sys.exit()

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
