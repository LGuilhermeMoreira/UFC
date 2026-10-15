from fastapi import APIRouter, Depends, Query
from typing import Optional, Dict, Any
from sqlalchemy.orm import Session
from app.models.all_models import Ingrediente as ingrediente_model
from app.dto.ingrediente import IngredienteCreate, IngredienteUpdate
from app.config.database import get_db
from app.controller.ingrediente import IngredienteController

router_ingrediente = APIRouter()

@router_ingrediente.post("/", response_model=ingrediente_model,status_code=201)
def create_ingrediente(ingrediente_data: IngredienteCreate, db: Session = Depends(get_db)):
    return IngredienteController.create_ingrediente(ingrediente_data, db)

@router_ingrediente.get("/", response_model=Dict[str, Any],status_code=200)
def list_ingredientes(
    page: int = Query(1, description="Número da página", ge=1),
    limit: int = Query(10, description="Número de itens por página", ge=1),
    nome: Optional[str] = Query(None, description="Filtrar por nome"),
    estoque: Optional[bool] = Query(None, description="Filtrar por estoque"),
    quantidade_estoque_min: Optional[float] = Query(None, description="Filtrar por quantidade mínima em estoque"),
    quantidade_estoque_max: Optional[float] = Query(None, description="Filtrar por quantidade máxima em estoque"),
    peso: Optional[float] = Query(None, description="Filtrar por peso"),
    db: Session = Depends(get_db)
):
    return IngredienteController.list_ingredientes(
        db,
        page=page,
        limit=limit,
        nome=nome,
        estoque=estoque,
        quantidade_estoque_min=quantidade_estoque_min,
        quantidade_estoque_max=quantidade_estoque_max,
        peso=peso
    )


@router_ingrediente.get("/{ingrediente_id}", response_model=ingrediente_model, status_code=200)
def get_ingrediente(ingrediente_id: int, db: Session = Depends(get_db)):
    return IngredienteController.get_ingrediente(ingrediente_id, db)

@router_ingrediente.put("/{ingrediente_id}", response_model=ingrediente_model, status_code=200)
def update_ingrediente(ingrediente_id: int, ingrediente_data: IngredienteUpdate, db: Session = Depends(get_db)):
    return IngredienteController.update_ingrediente(ingrediente_id, ingrediente_data, db)

@router_ingrediente.delete("/{ingrediente_id}", status_code=204)
def delete_ingrediente(ingrediente_id: int, db: Session = Depends(get_db)):
    return {"ok": IngredienteController.delete_ingrediente(ingrediente_id, db)}

@router_ingrediente.get("/num/", status_code=200)
def get_num_ingredientes(db: Session = Depends(get_db)):
    return IngredienteController.num_ingrediente(db)