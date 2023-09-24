import uvicorn
from main import app
from utils import config as uc

if __name__ == "__main__":
    uvicorn.run(app, host=uc.BACKEND_IP, port=uc.BACKEND_PORT)
