from typing import List, Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Comanda as comanda_model
from fastapi import HTTPException
from sqlmodel import select, and_, func
from app.dto.comanda import ComandaCreate, ComandaUpdate
import math
import logging

logger = logging.getLogger("api_logger")


class ComandaController:
    @staticmethod
    def create_comanda(comanda_data: ComandaCreate, db: Session) -> comanda_model:
        try:
            db_comanda = comanda_model(**comanda_data.model_dump())
            db.add(db_comanda)
            db.commit()
            db.refresh(db_comanda)
            logger.info(f"Comanda criada com sucesso. ID: {db_comanda.id}")
            return db_comanda
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao criar comanda: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def list_comandas(
        db: Session,
        page: int = 1,
        limit: int = 10,
        id_cliente: Optional[int] = None,
        id_mesa: Optional[int] = None,
        status: Optional[str] = None,
    ) -> Dict[str, Any]:
        try:
            offset = (page - 1) * limit
            query = select(comanda_model)
            
            filters = []
            if id_cliente:
              filters.append(comanda_model.id_cliente == id_cliente)
            if id_mesa:
                filters.append(comanda_model.id_mesa == id_mesa)
            if status:
                filters.append(comanda_model.status.ilike(f"%{status}%"))

            if filters:
                query = query.where(and_(*filters))
        
            comandas = db.execute(query.offset(offset).limit(limit)).scalars().all()
            
            total_query = select(func.count(comanda_model.id))
            if filters:
                total_query = total_query.where(and_(*filters))
            
            total = db.execute(total_query).scalar()
            total_pages = math.ceil(total / limit)
            logger.info(f"Listagem de comandas realizada. Filtros: Cliente ID={id_cliente}, Mesa ID={id_mesa}, Status={status}, Página={page}, Limite={limit}")
            return {
                "data": comandas,
                "pagination": {
                    "total": total,
                    "currentPage": page,
                    "totalPages": total_pages,
                    "totalItemsPerPage": limit
                },
            }
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar comandas: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def get_comanda(comanda_id: int, db: Session) -> comanda_model:
        try:
            comanda = db.get(comanda_model, comanda_id)
            if not comanda:
                logger.warning(f"Comanda nao encontrada. ID: {comanda_id}")
                raise HTTPException(status_code=404, detail="Comanda not found")
            logger.info(f"Comanda obtida com sucesso. ID: {comanda_id}")
            return comanda
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao obter comanda: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def update_comanda(comanda_id: int, comanda_data: ComandaUpdate, db: Session) -> comanda_model:
        try:
            db_comanda = db.get(comanda_model, comanda_id)
            if not db_comanda:
                logger.warning(f"Comanda nao encontrada para atualizacao. ID: {comanda_id}")
                raise HTTPException(status_code=404, detail="Comanda not found")
            for key, value in comanda_data.model_dump(exclude_unset=True).items():
                setattr(db_comanda, key, value)
            db.add(db_comanda)
            db.commit()
            db.refresh(db_comanda)
            logger.info(f"Comanda atualizada com sucesso. ID: {comanda_id}")
            return db_comanda
        except HTTPException as e:
             raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao atualizar comanda. ID: {comanda_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def delete_comanda(comanda_id: int, db: Session) -> bool:
        try:
            db_comanda = db.get(comanda_model, comanda_id)
            if not db_comanda:
                logger.warning(f"Comanda nao encontrada para remocao. ID: {comanda_id}")
                raise HTTPException(status_code=404, detail="Comanda not found")
            db.delete(db_comanda)
            db.commit()
            logger.info(f"Comanda removida com sucesso. ID: {comanda_id}")
            return True
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao remover comanda. ID: {comanda_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def num_comanda(db: Session) -> Dict[str,int]:
        try:
            num = db.query(func.count(comanda_model.id)).scalar()
            logger.info(f"Quantidade de comandas: {num}")
            return {"quantidade":num}
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao pegar a quantidade de comandas: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            if num is None:
                return {"quantidade" : 0} 
            return {"quantidade": num}
    
    @staticmethod
    def listar_comandas_por_cliente(db: Session, cliente_id: int) -> List[comanda_model]:
        try:
            statement = select(comanda_model).where(comanda_model.id_cliente == cliente_id)
            comandas = db.execute(statement).scalars().all()
            logger.info(f"Listagem de comandas por cliente realizada. ID do cliente: {cliente_id}")
            return comandas
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar comandas por cliente: {e}")
            raise HTTPException(status_code=500, detail=str(e))
    
    @staticmethod
    def listar_comandas_por_ano(db: Session, ano: int) -> List[comanda_model]:
        try:
            statement = select(comanda_model).where(func.extract('year', comanda_model.data_hora_abertura) == ano)
            comandas = db.execute(statement).scalars().all()
            logger.info(f"Listagem de comandas por ano realizada. Ano: {ano}")
            return comandas
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar comandas por ano: {e}")
            raise HTTPException(status_code=500, detail=str(e))