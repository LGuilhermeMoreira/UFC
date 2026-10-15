from typing import Optional, Dict, Any, List
from sqlalchemy.orm import Session
from app.models.all_models import Mesa as mesa_model
from app.dto.mesa import MesaCreate, MesaUpdate
from fastapi import HTTPException
from sqlmodel import select, and_, func
import math
import logging

logger = logging.getLogger("api_logger")


class MesaController:
    @staticmethod
    def create_mesa(mesa_data: MesaCreate, db: Session) -> mesa_model:
        try:
            db_mesa = mesa_model(**mesa_data.model_dump())
            db.add(db_mesa)
            db.commit()
            db.refresh(db_mesa)
            logger.info(f"Mesa criada com sucesso. ID: {db_mesa.id}")
            return db_mesa
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao criar mesa: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def list_mesas(
        db: Session,
        page: int = 1,
        limit: int = 10,
        numero_mesa: Optional[int] = None,
        capacidade: Optional[int] = None,
        ocupada: Optional[bool] = None,
        numero_pessoas: Optional[int] = None
    ) -> Dict[str, Any]:
        try:
            offset = (page - 1) * limit
            query = select(mesa_model)
            
            filters = []
            if numero_mesa:
              filters.append(mesa_model.numero_mesa == numero_mesa)
            if capacidade:
               filters.append(mesa_model.capacidade == capacidade)
            if ocupada is not None:
              filters.append(mesa_model.ocupada == ocupada)
            if numero_pessoas:
                 filters.append(mesa_model.numero_pessoas == numero_pessoas)
            
            if filters:
                query = query.where(and_(*filters))
        
            mesas = db.execute(query.offset(offset).limit(limit)).scalars().all()
            
            total_query = select(func.count(mesa_model.id))
            if filters:
                 total_query = total_query.where(and_(*filters))
            
            total = db.execute(total_query).scalar()
            total_pages = math.ceil(total / limit)
            logger.info(f"Listagem de mesas realizada. Filtros: Numero Mesa={numero_mesa}, Capacidade={capacidade}, Ocupada={ocupada}, Numero Pessoas={numero_pessoas}, Pagina={page}, Limite={limit}")
            return {
                "data": mesas,
                "pagination": {
                    "total": total,
                    "currentPage": page,
                    "totalPages": total_pages,
                    "totalItemsPerPage": limit
                },
            }
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar mesas: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def get_mesa(mesa_id: int, db: Session) -> mesa_model:
        try:
            mesa = db.get(mesa_model, mesa_id)
            if not mesa:
                logger.warning(f"Mesa nao encontrada. ID: {mesa_id}")
                raise HTTPException(status_code=404, detail="Mesa not found")
            logger.info(f"Mesa obtida com sucesso. ID: {mesa_id}")
            return mesa
        except HTTPException as e:
           raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao obter mesa: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def update_mesa(mesa_id: int, mesa_data: MesaUpdate, db: Session) -> mesa_model:
        try:
            db_mesa = db.get(mesa_model, mesa_id)
            if not db_mesa:
                logger.warning(f"Mesa nao encontrada para atualizacao. ID: {mesa_id}")
                raise HTTPException(status_code=404, detail="Mesa not found")
            for key, value in mesa_data.model_dump(exclude_unset=True).items():
                setattr(db_mesa, key, value)
            db.add(db_mesa)
            db.commit()
            db.refresh(db_mesa)
            logger.info(f"Mesa atualizada com sucesso. ID: {mesa_id}")
            return db_mesa
        except HTTPException as e:
           raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao atualizar mesa. ID: {mesa_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def delete_mesa(mesa_id: int, db: Session) -> bool:
        try:
            db_mesa = db.get(mesa_model, mesa_id)
            if not db_mesa:
                logger.warning(f"Mesa nao encontrada para remocao. ID: {mesa_id}")
                raise HTTPException(status_code=404, detail="Mesa not found")
            db.delete(db_mesa)
            db.commit()
            logger.info(f"Mesa removida com sucesso. ID: {mesa_id}")
            return True
        except HTTPException as e:
           raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao remover mesa. ID: {mesa_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))
    
    @staticmethod
    def num_mesa(db: Session) -> Dict[str,int]:
        try:
            num = db.query(func.count(mesa_model.id)).scalar()
            logger.info(f"Quantidade de mesas: {num}")
            return {"quantidade":num}
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao pegar a quantidade de mesas: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            if num is None:
                return {"quantidade" : 0} 
            return {"quantidade": num}
    
    @staticmethod
    def get_clientes_e_comandas_da_mesa(mesa_id: int, db: Session) -> List[Dict[str, Any]]:
        try:
            mesa = db.get(mesa_model, mesa_id)
            if not mesa:
                logger.warning(f"Mesa não encontrada. ID: {mesa_id}")
                raise HTTPException(status_code=404, detail="Mesa not found")

            clientes_e_comandas = []
            for comanda in mesa.comandas:
                if comanda.cliente:
                    clientes_e_comandas.append({
                        "id_cliente": comanda.cliente.id,
                        "nome_cliente": comanda.cliente.nome,
                        "id_comanda": comanda.id
                    })

            logger.info(f"Clientes e comandas da mesa obtidos com sucesso. ID da mesa: {mesa_id}")
            return clientes_e_comandas

        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao obter clientes e comandas da mesa. ID da mesa: {mesa_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))