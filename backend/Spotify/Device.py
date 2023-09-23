from typing import Optional

class SpotifyDevice:
    def __init__(self, *args, **kwargs) -> None:
        self.id = kwargs.get("id")
        self.name = kwargs.get("name")
        self.type = kwargs.get("type")

    def get_data(self) -> Optional[dict]:
        return (
            {"id": self.id, "name": self.name, "type": self.type}
            if self.validate_device()
            else None
        )

    def set_data(self, data: dict) -> None:
        self.id = data.get("id")
        self.name = data.get("name")
        self.type = data.get("type")

    def validate_device(self) -> bool:
        return self.id is not None
