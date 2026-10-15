from fastapi import APIRouter, Depends, Query
from sqlalchemy.orm import Session
from app.config.database import get_db
from app.controller.prato_ingrediente import PratoIngredienteController
from typing import List
from app.models.all_models import Ingrediente as ingrediente_model

router_prato_ingrediente = APIRouter()

@router_prato_ingrediente.post("/prato/{prato_id}/ingrediente/{ingrediente_id}", status_code=201)
def add_ingrediente_to_prato(
    prato_id: int,
    ingrediente_id: int,
    quantidade_utilizada: float = Query(1.0, description="Quantidade utilizada do ingrediente"),
    db: Session = Depends(get_db)
):
  return {"ok": PratoIngredienteController.add_ingrediente_to_prato(prato_id, ingrediente_id, quantidade_utilizada, db)}


@router_prato_ingrediente.delete("/prato/{prato_id}/ingrediente/{ingrediente_id}", status_code=204)
def remove_ingrediente_from_prato(prato_id: int, ingrediente_id: int, db: Session = Depends(get_db)):
    return {"ok": PratoIngredienteController.remove_ingrediente_from_prato(prato_id, ingrediente_id, db)}

@router_prato_ingrediente.get("/prato/{prato_id}/ingredientes", response_model=List[ingrediente_model], status_code=200)
def get_ingredientes_from_prato(prato_id: int, db: Session = Depends(get_db)):
    return PratoIngredienteController.get_ingredientes_from_prato(prato_id, db)