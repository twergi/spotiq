from fastapi import FastAPI
import os, requests, json, base64
from urllib.parse import quote_plus
from fastapi.responses import RedirectResponse
from typing import Optional
import datetime as dt

app = FastAPI()

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SECRET = os.environ.get("CLIENT_SECRET")
assert CLIENT_ID is not None, ""
assert CLIENT_SECRET is not None


API_URL = "https://api.spotify.com/v1/"

class SpotifyToken:

    token_url = "https://accounts.spotify.com/api/token"
    login_url = "https://accounts.spotify.com/authorize"

    def __init__(self, *args, **kwargs):
        self.access_token = kwargs.get("access_token")
        self.token_type = kwargs.get("token_type")
        self.refresh_token = kwargs.get("refresh_token")
        self.expires = kwargs.get("expires")

    @staticmethod
    def _calculate_expires(expires_in: int) -> dt.datetime:
        return dt.datetime.now() + dt.timedelta(seconds=expires_in)

    def create_auth_url(self) -> Optional[str]:
        response_type = "code"
        redirect_uri = "http://localhost:8000/callback/"
        scopes = (
            "user-read-playback-state",
            "user-read-currently-playing",
        )
        url = requests.Request(
            method="GET",
            url=self.login_url,
            params={
                "response_type": response_type,
                "redirect_uri": redirect_uri,
                "client_id": CLIENT_ID,
                "scope": " ".join(scopes)
            }
        ).prepare().url
        return url

    @property
    def expired(self):
        if self.expires is None:
            raise Exception("Token expiry time has not been set")
        
        return self.expires > dt.datetime.now()

    def get_current_queue(self):
        url = API_URL + "me/player/queue"

        response: dict = requests.get(
            url=url,
            headers={
                "Authorization": f"Bearer {self.access_token}"
            }
        ).json()
        return response

    def refresh(self, code: str):
        grant_type = "authorization_code"
        redirect_uri = "http://localhost:8000/callback/"

        response: dict = requests.post(
            url=self.token_url,
            data={
                "client_id": CLIENT_ID,
                "client_secret": CLIENT_SECRET,
                "redirect_uri": redirect_uri,
                "grant_type": grant_type,
                "code": code
            }
        ).json()
        error = response.get("error")

        if error:
            print(error)

        self.access_token = response.get("access_token")
        self.token_type = response.get("token_type")
        self.refresh_token = response.get("refresh_token")
        self.expires = token._calculate_expires(response.get("expires_in"))


token = SpotifyToken()


@app.get("/")
async def root():
    return {
        "token": {
            "access_token": token.access_token,
            "token_type": token.token_type,
            "refresh_token": token.refresh_token,
            "expires": token.expires
        }
    }

@app.get("/login/", response_class=RedirectResponse, status_code=302)
async def login():
    url = token.create_auth_url()
    return url

@app.get("/callback/", response_class=RedirectResponse, status_code=302)
async def spotify_callback(code: str, state: Optional[str] = None):
    token.refresh(code)
    return "http://localhost:8000/"

@app.get("/current_queue/")
async def get_current_queue():
    return token.get_current_queue()
