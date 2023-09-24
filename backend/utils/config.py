import os


DOCKER_RUN = True if os.environ.get("DOCKER_RUN") == "True" else False


BACKEND_IP = os.environ.get("BACKEND_IP") if DOCKER_RUN else "localhost"
BACKEND_PORT = int(os.environ.get("BACKEND_PORT", 0)) if DOCKER_RUN else 8000
BACKEND_PORT_FORWARD = int(os.environ.get("BACKEND_PORT_FORWARD", 0)) if DOCKER_RUN else 8000


FRONTEND_IP = os.environ.get("FRONTEND_IP") if DOCKER_RUN else "localhost"
FRONTEND_PORT = int(os.environ.get("FRONTEND_PORT_FORWARD", 0)) if DOCKER_RUN else 3000
FRONTEND_URL = f"http://{FRONTEND_IP}:{FRONTEND_PORT}"


REDIRECT_URI = f"http://{FRONTEND_IP}:{BACKEND_PORT_FORWARD}"
CLIENT_ID = os.environ.get("CLIENT_ID")
CLIENT_SECRET = os.environ.get("CLIENT_SECRET")


assert CLIENT_ID is not None, "CLIENT_ID must be specified in env"
assert CLIENT_SECRET is not None, "CLIENT_SECRET must be specified in env"
assert BACKEND_IP is not None, "BACKEND_IP must be specified in env"
assert BACKEND_PORT is not 0, "BACKEND_PORT must be specified in env"
assert BACKEND_PORT_FORWARD is not 0, "BACKEND_PORT must be specified in env"
assert FRONTEND_IP is not None, "FRONTEND_IP must be specified in env"
assert FRONTEND_PORT is not 0, "FRONTEND_PORT_FORWARD must be specified in env"