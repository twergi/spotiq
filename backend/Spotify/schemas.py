from pydantic import BaseModel

class Device(BaseModel):
    id: str
    name: str
    type: str


class Song(BaseModel):
    uri: str
