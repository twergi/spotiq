from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app import api_router

app = FastAPI()

origins = [
    "http://localhost:3000",
    "localhost:3000",
    "http://localhost",
    "localhost",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(api_router)
