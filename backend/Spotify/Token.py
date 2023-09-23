from fastapi import HTTPException, status
from typing import Optional
from utils.config import BACKEND_URL
import datetime as dt
import requests
from utils.config import CLIENT_ID, CLIENT_SECRET


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
            "user-modify-playback-state",
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



