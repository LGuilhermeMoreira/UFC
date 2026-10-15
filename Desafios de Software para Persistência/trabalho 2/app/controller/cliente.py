from typing import Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Cliente as cliente_model
from app.models.all_models import Comanda as comanda_model
from app.models.all_models import ComandaPratoLink as comandapratolink_model
from app.dto.cliente import ClienteCreate, ClienteUpdate
from fastapi import HTTPException
from sqlmodel import select, and_
from sqlalchemy import func
import math
import logging

logger = logging.getLogger("api_logger")

class ClienteController:
    @staticmethod
    def create_cliente(cliente_data: ClienteCreate, db: Session) -> cliente_model:
        try:
            db_cliente = cliente_model(**cliente_data.model_dump())
            db.add(db_cliente)
        except Exception as e:           
            db.rollback()
            logger.error(f"Erro ao criar cliente: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            db.commit()
            db.refresh(db_cliente)
            logger.info(f"Cliente criado com sucesso. ID: {db_cliente.id}")
            return db_cliente      

    @staticmethod
    def list_clientes(
        db: Session,
        page: int = 1,
        limit: int = 10,
        nome: Optional[str] = None,
        email: Optional[str] = None,
        cpf: Optional[str] = None,
    ) -> Dict[str, Any]:
        try:
            offset = (page - 1) * limit
            query = select(cliente_model)
            
            filters = []
            if nome:
              filters.append(cliente_model.nome.ilike(f"%{nome}%"))
            if email:
               filters.append(cliente_model.email.ilike(f"%{email}%"))
            if cpf:
                filters.append(cliente_model.cpf.ilike(f"%{cpf}%"))
          
            if filters:
                query = query.where(and_(*filters))
        
            clientes = db.execute(query.offset(offset).limit(limit)).scalars().all()
            
            total_query = select(func.count(cliente_model.id))
            if filters:
                total_query = total_query.where(and_(*filters))

            total = db.execute(total_query).scalar()
            total_pages = math.ceil(total / limit)
            logger.info(f"Listagem de clientes realizada. Filtros: Nome={nome}, Email={email}, CPF={cpf}, Pagina={page}, Limite={limit}")
            return {
                "data": clientes,
                "pagination": {
                    "total": total,
                    "currentPage": page,
                    "totalPages": total_pages,
                    "totalItemsPerPage": limit
                },
            }
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao listar clientes: {e}")
            raise HTTPException(status_code=500, detail=str(e))
    
    @staticmethod
    def get_cliente(cliente_id: int, db: Session) -> cliente_model:
        try:
            cliente = db.get(cliente_model, cliente_id)
            if not cliente:
                logger.warning(f"Cliente nao encontrado. ID: {cliente_id}")
                raise HTTPException(status_code=404, detail="Cliente not found")
            logger.info(f"Cliente obtido com sucesso. ID: {cliente_id}")
            return cliente
        except HTTPException as e:
           raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao obter cliente: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        
    @staticmethod
    def update_cliente(cliente_id: int, cliente_data: ClienteUpdate, db: Session) -> cliente_model:
        try:
            db_cliente = db.get(cliente_model, cliente_id)
            if not db_cliente:
                logger.warning(f"Cliente nao encontrado para atualizacao. ID: {cliente_id}")
                raise HTTPException(status_code=404, detail="Cliente not found")
            for key, value in cliente_data.model_dump(exclude_unset=True).items():
                setattr(db_cliente, key, value)
            db.add(db_cliente)
            db.commit()
            db.refresh(db_cliente)
            logger.info(f"Cliente atualizado com sucesso. ID: {cliente_id}")
            return db_cliente
        except HTTPException as e:
            raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao atualizar cliente. ID: {cliente_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def delete_cliente(cliente_id: int, db: Session) -> bool:
        try:
           db_cliente = db.get(cliente_model, cliente_id)
           if not db_cliente:
              logger.warning(f"Cliente nao encontrado para remocao. ID: {cliente_id}")
              raise HTTPException(status_code=404, detail="Cliente not found")
           db.delete(db_cliente)
           db.commit()
           logger.info(f"Cliente removido com sucesso. ID: {cliente_id}")
           return True
        except HTTPException as e:
           raise
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao remover cliente. ID: {cliente_id}, Erro: {e}")
            raise HTTPException(status_code=500, detail=str(e))

    @staticmethod
    def num_cliente(db: Session) -> Dict[str, int]:
        num: Optional[int] = None 
        try:
            num = db.query(func.count(cliente_model.id)).scalar()
            logger.info(f"Quantidade de clientes: {num}")
            return {"quantidade": num}
        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao pegar a quantidade de clientes: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            if num is None:
                return {"quantidade" : 0} 
            return {"quantidade" : num}
        
    @staticmethod
    def add_prato_to_comanda(comanda_id: int, prato_id: int, quantidade: int, db: Session) -> bool:
        try:
            db_comanda = db.get(comanda_model, comanda_id)
            if not db_comanda:
                raise HTTPException(status_code=404, detail="Comanda not found")
            
            db_prato = db.get(prato_model, prato_id)
            if not db_prato:
              raise HTTPException(status_code=404, detail="Prato not found")
            
            existing_link = db.execute(select(comandapratolink_model).where(and_(
                comandapratolink_model.id_comanda == comanda_id,
                comandapratolink_model.id_prato == prato_id
            ))).scalar_one_or_none()

            if existing_link:
                existing_link.quantidade += quantidade
            else:
                db_comanda_prato_link = comandapratolink_model(id_comanda=comanda_id, id_prato=prato_id, quantidade=quantidade)
                db.add(db_comanda_prato_link)

        except Exception as e:
            db.rollback()
            logger.error(f"Erro ao adicionar prato à comanda: {e}")
            raise HTTPException(status_code=500, detail=str(e))
        finally:
            db.commit()
            logger.info(f"Prato {prato_id} adicionado à comanda {comanda_id} com sucesso.")
            return True

