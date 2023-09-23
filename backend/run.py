from fastapi import FastAPI, HTTPException, status
import os, requests
from fastapi.responses import RedirectResponse
from typing import Optional
import datetime as dt
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware


app = FastAPI()

origins = [
    "http://localhost:3000",
]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SECRET = os.environ.get("CLIENT_SECRET")
assert CLIENT_ID is not None, ""
assert CLIENT_SECRET is not None

BACKEND_URL = "http://localhost:8000"
FRONTEND_URL = "http://localhost:3000"

class SpotifyToken:
    token_url = "https://accounts.spotify.com/api/token"
    login_url = "https://accounts.spotify.com/authorize"

    def __init__(self, **kwargs):
        self.access_token = kwargs.get("access_token")
        self.token_type = kwargs.get("token_type")
        self.refresh_token = kwargs.get("refresh_token")
        self.expires_at = kwargs.get("expires_at")

    @staticmethod
    def _calculate_expires_at(expires_in: int) -> dt.datetime:
        return dt.datetime.now() + dt.timedelta(seconds=expires_in)

    def create_auth_url(self) -> Optional[str]:
        response_type = "code"
        redirect_uri = f"{BACKEND_URL}/callback/"
        scopes = (
            "user-read-playback-state",
            "user-read-currently-playing",
        )
        url = (
            requests.Request(
                method="GET",
                url=self.login_url,
                params={
                    "client_id": CLIENT_ID,
                    "scope": " ".join(scopes),
                    "response_type": response_type,
                    "redirect_uri": redirect_uri,
                },
            )
            .prepare()
            .url
        )
        return url

    @property
    def expired(self) -> bool:
        if self.expires_at is None:
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Token expiry time has not been set",
            )

        return (self.expires_at - dt.timedelta(seconds=10)) < dt.datetime.now()

    def get_data(self):
        return (
            {
                "access_token": self.access_token,
                "token_type": self.token_type,
                "refresh_token": self.refresh_token,
                "expires_at": self.expires_at,
            }
            if self.validate_token()
            else None
        )

    def validate_token(self) -> bool:
        if not self.access_token:
            return False

        expired = self.expired

        if expired:
            self.refresh()
            expired = self.expired

        return not expired

    def refresh(self):
        grant_type = "refresh_token"

        response = requests.post(
            url=self.token_url,
            data={"grant_type": grant_type, "refresh_token": self.refresh_token},
        )

        if response.status_code != 200:
            raise HTTPException(
                status_code=response.status_code, detail=response.content.decode()
            )

        response = response.json()

        access_token = response.get("access_token")
        refresh_token = response.get("refresh_token")
        expires_in = response.get("expires_in")
        token_type = response.get("token_type")
        error = response.get("error")

        if error:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=error)

        expires_at = self._calculate_expires_at(expires_in) if expires_in else None

        self.access_token = access_token or self.access_token
        self.refresh_token = refresh_token or self.refresh_token
        self.expires_at = expires_at or self.expires_at
        self.token_type = token_type or self.token_type

    def _create_token_data(self, code: str):
        grant_type = "authorization_code"
        redirect_uri = f"{BACKEND_URL}/callback/"

        response: dict = requests.post(
            url=self.token_url,
            data={
                "client_id": CLIENT_ID,
                "client_secret": CLIENT_SECRET,
                "redirect_uri": redirect_uri,
                "grant_type": grant_type,
                "code": code,
            },
        ).json()
        error = response.get("error")

        if error:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=error)

        self.access_token = response.get("access_token")
        self.token_type = response.get("token_type")
        self.refresh_token = response.get("refresh_token")
        self.expires_at = self._calculate_expires_at(response.get("expires_in"))


class SpotifyDevice:
    def __init__(self, *args, **kwargs):
        self.id = kwargs.get("id")
        self.name = kwargs.get("name")
        self.type = kwargs.get("type")

    def get_data(self):
        return (
            {"id": self.id, "name": self.name, "type": self.type}
            if self.validate_device()
            else None
        )

    def set_data(self, data: dict):
        self.id = data.get("id")
        self.name = data.get("name")
        self.type = data.get("type")

    def validate_device(self):
        return self.id is not None


class SpotifyAPI:
    base_url = "https://api.spotify.com/v1"

    def __init__(self):
        self.token = SpotifyToken()
        self.device = SpotifyDevice()

    def clear_data(self):
        self.token = SpotifyToken()
        self.device = SpotifyDevice()

    def _make_auth_headers(self):
        return {"Authorization": f"Bearer {self.token.access_token}"}

    def get_data(self):
        return {"token": self._get_token_data(), "device": self._get_device_data()}

    def _get_token_data(self):
        return self.token.get_data()

    def create_token_data(self, code: str):
        self.token._create_token_data(code)

    def _get_device_data(self):
        return self.device.get_data()

    def get_devices(self):
        self.token.validate_token()

        response = requests.get(
            url=self.base_url + "/me/player/devices", headers=self._make_auth_headers()
        )
        if response.status_code != 200:
            raise HTTPException(status_code=response.status_code, detail=response.content.decode())

        return response.json()

    def set_device(self, device_info: dict):
        self.device.set_data(device_info)

    def get_current_queue(self):
        self.token.validate_token()

        url = self.base_url + "/me/player/queue"

        response: dict = requests.get(url=url, headers=self._make_auth_headers()).json()

        return response

    def get_currently_playing(self):
        self.token.validate_token()

        url = self.base_url + "/me/player/currently-playing"

        response: dict = requests.get(url=url, headers=self._make_auth_headers()).json()

        return response

    def search(self, query: str, **kwargs) -> dict:
        self.token.validate_token()

        url = f"{self.base_url}/search?q={query}&type=track&limit=10"

        response: dict = requests.get(url=url, headers=self._make_auth_headers())

        error = response.get("error")
        if error:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=error)

        return response["tracks"]


api_obj = SpotifyAPI()


@app.get("/")
async def root():
    return api_obj.get_data()


@app.get(
    "/login/",
    response_class=RedirectResponse,
    status_code=status.HTTP_307_TEMPORARY_REDIRECT,
)
async def login():
    return api_obj.token.create_auth_url()


@app.get(
    "/logout/",
    response_class=RedirectResponse,
    status_code=status.HTTP_307_TEMPORARY_REDIRECT,
)
async def login():
    api_obj.clear_data()
    return BACKEND_URL


@app.get(
    "/callback/",
    response_class=RedirectResponse,
    status_code=status.HTTP_307_TEMPORARY_REDIRECT,
)
async def spotify_callback(code: str, state: Optional[str] = None):
    api_obj.create_token_data(code)
    return FRONTEND_URL


@app.get("/current_queue/")
async def get_current_queue():
    return api_obj.get_current_queue()


@app.get("/currently_playing/")
async def get_currently_playing():
    return api_obj.get_currently_playing()


@app.get("/search")
async def search(q: str):
    return api_obj.search(query=q)


@app.get("/devices/")
async def get_devices():
    return api_obj.get_devices()


class Device(BaseModel):
    id: str
    name: str
    type: str


@app.post("/devices/")
async def set_device(device_info: Device):
    return api_obj.set_device(device_info)
