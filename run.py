from fastapi import FastAPI
import os, requests, json, base64

app = FastAPI()

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SECRET = os.environ.get("CLIENT_SECRET")

assert CLIENT_ID is not None
assert CLIENT_SECRET is not None

class SpotifyToken:
    def __init__(self, *args, **kwargs):
        self.access_token = kwargs["access_token"]
        self.token_type = kwargs["token_type"]

token = SpotifyToken(
    access_token=None,
    token_type=None
)

API_URL = "https://api.spotify.com/v1/"
TOKEN_URL = "https://accounts.spotify.com/api/token"

async def update_token():
    encoded_data = base64.b64encode(
        (CLIENT_ID+':'+CLIENT_SECRET).encode("ascii")
    )
    response = requests.post(
        TOKEN_URL,
        headers={
            "Content-Type": "application/x-www-form-urlencoded",
            "Authorization": f"Basic {encoded_data}"
        },
        data={"grant_type": "client_credentials"}
    )

    data = json.loads(
        response.content.decode()
    )
    print(data)
    token.access_token = data.get("access_token")
    token.token_type = data.get("token_type")

@app.get("/")
async def root():
    return {"message": "Hello World"}

@app.get("/update_token/")
async def get_update_token():
    await update_token()
    return None

@app.get("/current_queue/")
async def get_current_queue():
    response = requests.get(
        API_URL + "me/player/queue",
        headers={
            "Authorization": f"Bearer {token.access_token}"
        }
    )

    if response.status_code == 401:
        await update_token()
        response = requests.get(
            API_URL + "me/player/queue",
            headers={
                "Authorization": f"Bearer {token.access_token}"
            }
        )

    # print(token.access_token)
    return {"data": json.loads(response.content.decode())}
