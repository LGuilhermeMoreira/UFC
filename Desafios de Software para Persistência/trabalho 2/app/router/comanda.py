from fastapi import APIRouter, Depends, Query,Path
from typing import Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Comanda as comanda_model
from app.dto.comanda import ComandaCreate, ComandaUpdate
from app.config.database import get_db
from app.controller.comanda import ComandaController

router_comanda = APIRouter()

@router_comanda.post("/", response_model=comanda_model,status_code=201)
def create_comanda(comanda_data: ComandaCreate, db: Session = Depends(get_db),):
    return ComandaController.create_comanda(comanda_data, db)


@router_comanda.get("/", response_model=Dict[str, Any])
def list_comandas(
    page: int = Query(1, description="Número da página", ge=1),
    limit: int = Query(10, description="Número de itens por página", ge=1),
    id_cliente: Optional[int] = Query(None, description="Filtrar por ID do Cliente"),
    id_mesa: Optional[int] = Query(None, description="Filtrar por ID da Mesa"),
    status: Optional[str] = Query(None, description="Filtrar por Status"),
    db: Session = Depends(get_db)
):
    return ComandaController.list_comandas(db, page=page, limit=limit, id_cliente=id_cliente, id_mesa=id_mesa, status=status)


@router_comanda.get("/{comanda_id}", response_model=comanda_model,status_code=200)
def get_comanda(comanda_id: int, db: Session = Depends(get_db)):
    return ComandaController.get_comanda(comanda_id, db)

@router_comanda.put("/{comanda_id}", response_model=comanda_model,status_code=200)
def update_comanda(comanda_id: int, comanda_data: ComandaUpdate, db: Session = Depends(get_db)):
    return ComandaController.update_comanda(comanda_id, comanda_data, db)

@router_comanda.delete("/{comanda_id}",status_code=204)
def delete_comanda(comanda_id: int, db: Session = Depends(get_db)):
    return {"ok": ComandaController.delete_comanda(comanda_id, db)}

@router_comanda.get("/num/",status_code=200)
def get_num_comandas(db: Session = Depends(get_db)):
    return ComandaController.num_comanda(db)

@router_comanda.get("/comandas_por_cliente/{cliente_id}",status_code=200)
def comanda_por_cliente(db: Session = Depends(get_db),cliente_id: int = Path()):
    return ComandaController.listar_comandas_por_cliente(db,cliente_id)

@router_comanda.get("/comadas_por_ano/{ano}",status_code=200)
def comandas_por_ano(db: Session = Depends(get_db),ano: int = Path()):
    return ComandaController.listar_comandas_por_ano(db,ano)

