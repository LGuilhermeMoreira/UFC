from typing import List
from sqlalchemy.orm import Session
from app.models.all_models import Prato as prato_model
from app.models.all_models import Ingrediente as ingrediente_model
from app.models.all_models import PratoIngredienteLink as prato_ingrediente_link_model
from fastapi import HTTPException
from sqlmodel import select, and_
import logging

logger = logging.getLogger("api_logger")

class PratoIngredienteController:
    @staticmethod
    def add_ingrediente_to_prato(prato_id: int, ingrediente_id: int, quantidade_utilizada: float, db: Session) -> bool:
        try:
            db_prato = db.get(prato_model, prato_id)
            if not db_prato:
                raise HTTPException(status_code=404, detail="Prato not found")

            db_ingrediente = db.get(ingrediente_model, ingrediente_id)
            if not db_ingrediente:
                raise HTTPException(status_code=404, detail="Ingrediente not found")

            existing_link = db.execute(select(prato_ingrediente_link_model).where(and_(
                prato_ingrediente_link_model.id_prato == prato_id,
                prato_ingrediente_link_model.id_ingrediente == ingrediente_id
            ))).scalar_one_or_none()

            if existing_link:
                existing_link.quantidade_utilizada += quantidade_utilizada
            else:
               db_prato_ingrediente_link = prato_ingrediente_link_model(id_prato=prato_id, id_ingrediente=ingrediente_id, quantidade_utilizada=quantidade_utilizada)
               db.add(db_prato_ingrediente_link)

        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao adicionar ingrediente ao prato: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            db.commit()
            logger.info(f"Ingrediente {ingrediente_id} adicionado ao prato {prato_id} com sucesso.")
            return True

    @staticmethod
    def remove_ingrediente_from_prato(prato_id: int, ingrediente_id: int, db: Session) -> bool:
        try:
            db_prato_ingrediente_link = db.execute(
              select(prato_ingrediente_link_model).where(and_(
                  prato_ingrediente_link_model.id_prato == prato_id,
                  prato_ingrediente_link_model.id_ingrediente == ingrediente_id
                  )
              )
            ).scalar_one_or_none()
            if not db_prato_ingrediente_link:
                raise HTTPException(status_code=404, detail="Ingrediente not found in this plate")
            db.delete(db_prato_ingrediente_link)

        except Exception as e:
          db.rollback()
          logger.error(f"Erro ao remover o ingrediente do prato: {e}")
          raise HTTPException(status_code=500, detail=str(e))
        finally:
          db.commit()
          logger.info(f"Ingrediente {ingrediente_id} removido do prato {prato_id} com sucesso.")
          return True
    
    @staticmethod
    def get_ingredientes_from_prato(prato_id: int, db: Session) -> List[ingrediente_model]:
      try:
          db_prato = db.get(prato_model, prato_id)
          if not db_prato:
            raise HTTPException(status_code=404, detail="Prato not found")

          query = select(ingrediente_model).join(
            prato_ingrediente_link_model,
              and_(
                  prato_ingrediente_link_model.id_ingrediente == ingrediente_model.id,
                  prato_ingrediente_link_model.id_prato == prato_id
              )
          )

          ingredientes = db.execute(query).scalars().all()
      except Exception as e:
          db.rollback()
          logger.error(f"Erro ao listar ingredientes do prato: {e}")
          raise HTTPException(status_code=500, detail=str(e))
      finally:
          return ingredientes