from typing import Optional
from fastapi import APIRouter, status
from fastapi.responses import RedirectResponse
from Spotify import api_obj
from Spotify.schemas import Device, Song
from utils.config import FRONTEND_URL

router = APIRouter()


@router.get("/")
async def root():
    return api_obj.get_data()


@router.get(
    "/login/",
    response_class=RedirectResponse,
    status_code=status.HTTP_307_TEMPORARY_REDIRECT,
)
async def login():
    return api_obj.token.create_auth_url()


@router.get(
    "/logout/",
    response_class=RedirectResponse,
    status_code=status.HTTP_307_TEMPORARY_REDIRECT,
)
async def logout():
    api_obj.clear_data()
    return FRONTEND_URL


@router.get(
    "/callback/",
    response_class=RedirectResponse,
    status_code=status.HTTP_307_TEMPORARY_REDIRECT,
)
async def spotify_callback(code: str, state: Optional[str] = None):
    api_obj.create_token_data(code)
    return FRONTEND_URL


@router.get("/current_queue/")
async def get_current_queue():
    return api_obj.get_current_queue()




@router.post("/current_queue/")
async def get_current_queue(song_info: Song):
    return api_obj.add_to_queue(song_info.model_dump())


@router.get("/currently_playing/")
async def get_currently_playing():
    return api_obj.get_currently_playing()


@router.get("/search")
async def search(q: str):
    return api_obj.search(query=q)


@router.get("/devices/")
async def get_devices():
    return api_obj.get_devices()




@router.post("/devices/")
async def set_device(device_info: Device):
    return api_obj.set_device(device_info.model_dump())
