from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app import api_router
from utils import config as uc

app = FastAPI()

origins = (
    "http://localhost:3000",
    "localhost:3000",
    "http://localhost",
    "localhost",
    f"{uc.FRONTEND_IP}:{uc.FRONTEND_PORT}",
    f"http://{uc.FRONTEND_IP}:{uc.FRONTEND_PORT}",
    f"{uc.FRONTEND_IP}",
    f"http://{uc.FRONTEND_IP}",
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(api_router)
