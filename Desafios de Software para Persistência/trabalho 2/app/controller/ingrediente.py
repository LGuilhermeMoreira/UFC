from typing import Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Ingrediente as ingrediente_model
from app.dto.ingrediente import IngredienteCreate, IngredienteUpdate
from fastapi import HTTPException
from sqlmodel import select, and_, func
import math
import logging

logger = logging.getLogger("api_logger")

class IngredienteController:
    @staticmethod
    def create_ingrediente(ingrediente_data: IngredienteCreate, db: Session) -> ingrediente_model:
        try:
            db_ingrediente = ingrediente_model(**ingrediente_data.model_dump())
            db.add(db_ingrediente)
            db.commit()
            db.refresh(db_ingrediente)
            logger.info(f"Ingrediente criado com sucesso. ID: {db_ingrediente.id}")
            return db_ingrediente
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao criar ingrediente: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def list_ingredientes(
        db: Session,
        page: int = 1,
        limit: int = 10,
        nome: Optional[str] = None,
        estoque: Optional[bool] = None,
        quantidade_estoque_min: Optional[float] = None,
        quantidade_estoque_max: Optional[float] = None,
        peso: Optional[float] = None,
    ) -> Dict[str, Any]:
        try:
            offset = (page - 1) * limit
            query = select(ingrediente_model)
            
            filters = []
            if nome:
                filters.append(ingrediente_model.nome.ilike(f"%{nome}%"))
            if estoque is not None:
                filters.append(ingrediente_model.estoque == estoque)
            if quantidade_estoque_min is not None:
                filters.append(ingrediente_model.quantidade_estoque >= quantidade_estoque_min)
            if quantidade_estoque_max is not None:
                filters.append(ingrediente_model.quantidade_estoque <= quantidade_estoque_max)
            if peso is not None:
                filters.append(ingrediente_model.peso == peso)

            if filters:
                query = query.where(and_(*filters))
        
            ingredientes = db.execute(query.offset(offset).limit(limit)).scalars().all()
            
            total_query = select(func.count(ingrediente_model.id))
            if filters:
                total_query = total_query.where(and_(*filters))
            
            total = db.execute(total_query).scalar()
            total_pages = math.ceil(total / limit)
            logger.info(f"Listagem de ingredientes realizada. Filtros: Nome={nome}, Estoque={estoque}, Quantidade Min={quantidade_estoque_min}, Quantidade Max={quantidade_estoque_max}, Peso={peso}, Pagina={page}, Limite={limit}")
            return {
                "data": ingredientes,
                "pagination": {
                    "total": total,
                    "currentPage": page,
                    "totalPages": total_pages,
                    "totalItemsPerPage": limit
                },
            }
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar ingredientes: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def get_ingrediente(ingrediente_id: int, db: Session) -> ingrediente_model:
        try:
            ingrediente = db.get(ingrediente_model, ingrediente_id)
            if not ingrediente:
                logger.warning(f"Ingrediente nao encontrado. ID: {ingrediente_id}")
                raise HTTPException(status_code=404, detail="Ingrediente not found")
            logger.info(f"Ingrediente obtido com sucesso. ID: {ingrediente_id}")
            return ingrediente
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao obter ingrediente: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def update_ingrediente(ingrediente_id: int, ingrediente_data: IngredienteUpdate, db: Session) -> ingrediente_model:
        try:
            db_ingrediente = db.get(ingrediente_model, ingrediente_id)
            if not db_ingrediente:
                logger.warning(f"Ingrediente nao encontrado para atualizacao. ID: {ingrediente_id}")
                raise HTTPException(status_code=404, detail="Ingrediente not found")
            for key, value in ingrediente_data.model_dump(exclude_unset=True).items():
                setattr(db_ingrediente, key, value)
            db.add(db_ingrediente)
            db.commit()
            db.refresh(db_ingrediente)
            logger.info(f"Ingrediente atualizado com sucesso. ID: {ingrediente_id}")
            return db_ingrediente
        except HTTPException as e:
           raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao atualizar ingrediente. ID: {ingrediente_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def delete_ingrediente(ingrediente_id: int, db: Session) -> bool:
        try:
            db_ingrediente = db.get(ingrediente_model, ingrediente_id)
            if not db_ingrediente:
                logger.warning(f"Ingrediente nao encontrado para remocao. ID: {ingrediente_id}")
                raise HTTPException(status_code=404, detail="Ingrediente not found")
            db.delete(db_ingrediente)
            db.commit()
            logger.info(f"Ingrediente removido com sucesso. ID: {ingrediente_id}")
            return True
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao remover ingrediente. ID: {ingrediente_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))
    
    @staticmethod
    def num_ingrediente(db: Session) -> Dict[str,int]:
        try:
           num = db.query(func.count(ingrediente_model.id)).scalar()
           logger.info(f"Quantidade de ingredientes: {num}")
           return {"quantidade":num}
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao pegar a quantidade de ingredientes: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            if num is None:
                return {"quantidade" : 0} 
            return {"quantidade": num}