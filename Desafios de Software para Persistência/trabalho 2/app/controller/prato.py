from typing import Optional, Dict, Any, List
from sqlalchemy.orm import Session
from app.models.all_models import Prato as prato_model
from app.models.all_models import PratoIngredienteLink
from app.dto.prato import PratoCreate, PratoUpdate
from fastapi import HTTPException
from sqlmodel import select, and_, func
import math
import logging

logger = logging.getLogger("api_logger")

class PratoController:
    @staticmethod
    def create_prato(prato_data: PratoCreate, db: Session) -> prato_model:
        try:
            db_prato = prato_model(**prato_data.model_dump())
            db.add(db_prato)
            db.commit()
            db.refresh(db_prato)
            logger.info(f"Prato criado com sucesso. ID: {db_prato.id}")
            return db_prato
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao criar prato: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def list_pratos(
        db: Session,
        page: int = 1,
        limit: int = 10,
        nome: Optional[str] = None,
        preco_min: Optional[float] = None,
        preco_max: Optional[float] = None,
        disponivel: Optional[bool] = None,
    ) -> Dict[str, Any]:
        try:
            offset = (page - 1) * limit
            query = select(prato_model)
            
            filters = []
            if nome:
                filters.append(prato_model.nome.ilike(f"%{nome}%"))
            if preco_min is not None:
                filters.append(prato_model.preco >= preco_min)
            if preco_max is not None:
                filters.append(prato_model.preco <= preco_max)
            if disponivel is not None:
                filters.append(prato_model.disponivel == disponivel)
                
            if filters:
                query = query.where(and_(*filters))
        
            pratos = db.execute(query.offset(offset).limit(limit)).scalars().all()
            
            total_query = select(func.count(prato_model.id))
            if filters:
                total_query = total_query.where(and_(*filters))
            
            total = db.execute(total_query).scalar()
            total_pages = math.ceil(total / limit)
            logger.info(f"Listagem de pratos realizada. Filtros: Nome={nome}, Preco Min={preco_min}, Preco Max={preco_max}, Disponivel={disponivel}, Pagina={page}, Limite={limit}")
            return {
                "data": pratos,
                "pagination": {
                    "total": total,
                    "currentPage": page,
                    "totalPages": total_pages,
                    "totalItemsPerPage": limit
                },
            }
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar pratos: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def get_prato(prato_id: int, db: Session) -> prato_model:
        try:
            prato = db.get(prato_model, prato_id)
            if not prato:
                logger.warning(f"Prato nao encontrado. ID: {prato_id}")
                raise HTTPException(status_code=404, detail="Prato not found")
            logger.info(f"Prato obtido com sucesso. ID: {prato_id}")
            return prato
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao obter prato: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def update_prato(prato_id: int, prato_data: PratoUpdate, db: Session) -> prato_model:
        try:
            db_prato = db.get(prato_model, prato_id)
            if not db_prato:
                logger.warning(f"Prato nao encontrado para atualizacao. ID: {prato_id}")
                raise HTTPException(status_code=404, detail="Prato not found")
            for key, value in prato_data.model_dump(exclude_unset=True).items():
                setattr(db_prato, key, value)
            db.add(db_prato)
            db.commit()
            db.refresh(db_prato)
            logger.info(f"Prato atualizado com sucesso. ID: {prato_id}")
            return db_prato
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao atualizar prato. ID: {prato_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def delete_prato(prato_id: int, db: Session) -> bool:
        try:
            db_prato = db.get(prato_model, prato_id)
            if not db_prato:
                logger.warning(f"Prato nao encontrado para remocao. ID: {prato_id}")
                raise HTTPException(status_code=404, detail="Prato not found")
            db.delete(db_prato)
            db.commit()
            logger.info(f"Prato removido com sucesso. ID: {prato_id}")
            return True
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao remover prato. ID: {prato_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))
    
    @staticmethod
    def num_prato(db: Session) -> Dict[str,int]:
        try:
            num = db.query(func.count(prato_model.id)).scalar()
            logger.info(f"Quantidade de pratos: {num}")
            return {"quantidade":num}
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao pegar a quantidade de pratos: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            if num is None:
                return {"quantidade" : 0} 
            return {"quantidade": num}
    
    @staticmethod
    def listar_pratos_ordenados_por_preco(db: Session) -> List[prato_model]:
        try:
            statement = select(prato_model).order_by(prato_model.preco)
            pratos = db.execute(statement).scalars().all()
            logger.info("Listagem de pratos ordenada por preco realizada")
            return pratos
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar pratos ordenados por preco: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            if pratos is not None:
                return pratos
            return []
    

    