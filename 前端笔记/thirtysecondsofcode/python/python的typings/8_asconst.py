from typing import Final, TypeVar


MAX_SIZE: Final = 9000
MAX_SIZE += 1  # Error reported by type checker


class Connection:
    TIMEOUT: Final[int] = 10


class FastConnector(Connection):
    TIMEOUT = 1  # Error reported by type checker

