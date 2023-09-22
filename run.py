from fastapi import FastAPI, HTTPException, status
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
            raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Token expiry time has not been set")
        
        return self.expires - dt.timedelta(seconds=10) > dt.datetime.now()

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
        self.expires = self._calculate_expires(response.get("expires_in"))


class SpotifyAPI:
    base_url = "https://api.spotify.com/v1"
    
    def __init__(self):
        self.token_obj = SpotifyToken()

    def search(self, query: str) -> dict:
        print(f"{self.base_url}/search?q={query}&type=track&limit=10)")
        response: dict = requests.get(
            url=f"{self.base_url}/search?q={query}&type=track&limit=10",
            
            headers={
                "Authorization": f"Bearer {self.token_obj.access_token}"
            }
        )
        print(response)
        print(response.content)
        response=response.json()
        print(response)
        error = response.get("error")
        if error:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=error)
        return response["tracks"]


spotify_obj = SpotifyAPI()

@app.get("/")
async def root():
    return {
        "token": {
            "access_token": spotify_obj.token_obj.access_token,
            "token_type": spotify_obj.token_obj.token_type,
            "refresh_token": spotify_obj.token_obj.refresh_token,
            "expires": spotify_obj.token_obj.expires
        }
    }

@app.get("/login/", response_class=RedirectResponse, status_code=302)
async def login():
    url = spotify_obj.token_obj.create_auth_url()
    return url

@app.get("/callback/", response_class=RedirectResponse, status_code=302)
async def spotify_callback(code: str, state: Optional[str] = None):
    spotify_obj.token_obj.refresh(code)
    return "http://localhost:8000/"

@app.get("/current_queue/")
async def get_current_queue():
    return spotify_obj.token_obj.get_current_queue()

@app.get("/search")
async def search(q: str):
    return spotify_obj.search(query=q)
