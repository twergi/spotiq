from fastapi import APIRouter
from app.routing import router


api_router = APIRouter()

api_router.include_router(router)
