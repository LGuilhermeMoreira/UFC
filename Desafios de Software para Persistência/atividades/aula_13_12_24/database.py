
DATABASE_URL= 'sqlite:///./escola.db'
engine = create_database(DATABASE_URL,connect_args={
   "check_same_thread": False
})

Base.metadata.create_all(bind=engine)

SessionLocal = sessionmaker(autocommit=False,autoflush=False,bind=engine)

def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
