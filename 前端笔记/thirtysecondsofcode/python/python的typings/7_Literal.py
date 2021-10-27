from typing import Any, Literal


def validate_simple(data: Any) -> Literal[True]:  # always returns True
    ...


MODE = Literal['r', 'rb', 'w', 'wb']


def open_helper(file: str, mode: MODE) -> str:
    ...


open_helper('/some/path', 'r')  # Passes type check
open_helper('/other/path', 'typo')  # Error in type checker
