from fastapi import APIRouter, Depends, Query
from typing import Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Prato as prato_model
from app.dto.prato import PratoCreate, PratoUpdate
from app.config.database import get_db
from app.controller.prato import PratoController

router_prato = APIRouter()

@router_prato.post("/", response_model=prato_model, status_code=201)
def create_prato(prato_data: PratoCreate, db: Session = Depends(get_db)):
    return PratoController.create_prato(prato_data, db)

@router_prato.get("/", response_model=Dict[str, Any],status_code=200)
def list_pratos(
    page: int = Query(1, description="Número da página", ge=1),
    limit: int = Query(10, description="Número de itens por página", ge=1),
    nome: Optional[str] = Query(None, description="Filtrar por nome"),
    preco_min: Optional[float] = Query(None, description="Filtrar por preço mínimo"),
    preco_max: Optional[float] = Query(None, description="Filtrar por preço máximo"),
    disponivel: Optional[bool] = Query(None, description="Filtrar por disponibilidade"),
    db: Session = Depends(get_db)
):
    return PratoController.list_pratos(
        db,
        page=page,
        limit=limit,
        nome=nome,
        preco_min=preco_min,
        preco_max=preco_max,
        disponivel=disponivel,
    )

@router_prato.get("/{prato_id}", response_model=prato_model,status_code=200)
def get_prato(prato_id: int, db: Session = Depends(get_db)):
    return PratoController.get_prato(prato_id, db)

@router_prato.put("/{prato_id}", response_model=prato_model,status_code=200)
def update_prato(prato_id: int, prato_data: PratoUpdate, db: Session = Depends(get_db)):
    return PratoController.update_prato(prato_id, prato_data, db)

@router_prato.delete("/{prato_id}",status_code=204)
def delete_prato(prato_id: int, db: Session = Depends(get_db)):
    return {"ok": PratoController.delete_prato(prato_id, db)}

@router_prato.get("/num/",status_code=200)
def get_num_pratos(db: Session = Depends(get_db)):
    return PratoController.num_prato(db)

@router_prato.get("/ordenado_por_preco/",status_code=200)
def get_pratos_ordenados_por_preco(db: Session = Depends(get_db)):
    return PratoController.listar_pratos_ordenados_por_preco(db)
