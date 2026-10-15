from fastapi import APIRouter, Depends, Query, Path
from typing import Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Mesa as mesa_model
from app.dto.mesa import MesaCreate, MesaUpdate
from app.config.database import get_db
from app.controller.mesa import MesaController

router_mesa = APIRouter()

@router_mesa.post("/", response_model=mesa_model,status_code=201)
def create_mesa(mesa_data: MesaCreate, db: Session = Depends(get_db)):
    return MesaController.create_mesa(mesa_data, db)

@router_mesa.get("/", response_model=Dict[str, Any],status_code=200)
def list_mesas(
    page: int = Query(1, description="Número da página", ge=1),
    limit: int = Query(10, description="Número de itens por página", ge=1),
    numero_mesa: Optional[int] = Query(None, description="Filtrar por número da mesa"),
    capacidade: Optional[int] = Query(None, description="Filtrar por capacidade"),
    ocupada: Optional[bool] = Query(None, description="Filtrar por mesas ocupadas ou não"),
    numero_pessoas: Optional[int] = Query(None, description="Filtrar por número de pessoas"),
    db: Session = Depends(get_db)
):
    return MesaController.list_mesas(
        db,
        page=page,
        limit=limit,
        numero_mesa=numero_mesa,
        capacidade=capacidade,
        ocupada=ocupada,
        numero_pessoas=numero_pessoas
    )


@router_mesa.get("/{mesa_id}", response_model=mesa_model, status_code=200)
def get_mesa(mesa_id: int, db: Session = Depends(get_db)):
    return MesaController.get_mesa(mesa_id, db)

@router_mesa.put("/{mesa_id}", response_model=mesa_model,status_code=200)
def update_mesa(mesa_id: int, mesa_data: MesaUpdate, db: Session = Depends(get_db)):
    return MesaController.update_mesa(mesa_id, mesa_data, db)

@router_mesa.delete("/{mesa_id}",status_code=204)
def delete_mesa(mesa_id: int, db: Session = Depends(get_db)):
    return {"ok": MesaController.delete_mesa(mesa_id, db)}

@router_mesa.get("/num/",status_code=200)
def get_num_mesas(db: Session = Depends(get_db)):
    return MesaController.num_mesa(db)

@router_mesa.get("/cliente_e_comanda/{mesa_id}",status_code=200)
def get_cliente_and_comanda(db : Session = Depends(get_db),mesa_id :int =  Path()):
    return MesaController.get_clientes_e_comandas_da_mesa(mesa_id,db)