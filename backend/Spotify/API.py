from fastapi import HTTPException, status
from Spotify.Token import SpotifyToken
from Spotify.Device import SpotifyDevice
import requests


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
            raise HTTPException(
                status_code=response.status_code, detail=response.content.decode()
            )

        return response.json()

    def set_device(self, device_info: dict):
        self.device.set_data(device_info)

    def get_current_queue(self):
        self.token.validate_token()

        url = self.base_url + "/me/player/queue"

        response: dict = requests.get(url=url, headers=self._make_auth_headers()).json()

        return response

    def add_to_queue(self, song_info: dict) -> None:
        song_uri = song_info.get("uri")

        self.token.validate_token()

        url = "/me/player/queue"

        response = requests.post(
            url=self.base_url + url,
            headers=self._make_auth_headers(),
            params={"uri": song_uri, "device_id": self.device.id},
        )

        if response.status_code != 204:
            raise HTTPException(
                status_code=response.status_code, detail=response.content.decode()
            )

    def get_currently_playing(self):
        self.token.validate_token()

        url = self.base_url + "/me/player/currently-playing"

        response: dict = requests.get(url=url, headers=self._make_auth_headers()).json()

        return response

    def search(self, query: str, **kwargs) -> dict:
        self.token.validate_token()

        url = f"{self.base_url}/search?q={query}&type=track&limit=10"

        response: dict = requests.get(url=url, headers=self._make_auth_headers()).json()

        error = response.get("error")
        if error:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail=error)

        return response["tracks"]
