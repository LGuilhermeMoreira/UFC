
class Aluno(Base):
    __tablename__ = 'alunos'
    id = Column(Intenger,primary_key=True)
    name = Column(String,nullable=False)
    email = Column(String,unique=True,nullabel=False)

    incricoes = relationship('Incricao',back_populates='alunos')

class Curso(Base):
    __tablename__ = 'cursos'
    id = Column(Intenger,primary_key=True)
    name = Column(String,nullable=False)
    descricao - Column(String)

    incricoes = relationship('Incricao',back_populates='cursos')

class Incricao(Base):
    __tablename__ = 'inscricoes'
    id = Column(Integer,primary_key = True)
    aluno_id = Column(Integer,ForeingKey('alunos.id'))
    curso_id = Column(Integer,ForeingKey('cursos.id'))

    aluno = relationship('Aluno',back_populate='inscricoes')
    curso = relationship('Curso',back_populate='inscricoes')
