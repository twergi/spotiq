import os


BACKEND_URL = "http://localhost:8000"
FRONTEND_URL = "http://localhost:3000"

CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SECRET = os.environ.get("CLIENT_SECRET")
assert CLIENT_ID is not None, "CLIENT_ID must be specified in env"
assert CLIENT_SECRET is not None, "CLIENT_SECRET must be specified in env"
