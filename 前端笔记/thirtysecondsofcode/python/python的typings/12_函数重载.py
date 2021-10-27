from typing import Any, Optional, Tuple, Union, overload


@overload
def process(response: None) -> None:
    ...


@overload
def process(response: int) -> Tuple[int, str]:
    ...


@overload
def process(response: bytes) -> str:
    ...


def process(response: Any) -> Any:
    return response


####################################################

# Overload *variants* for 'mouse_event'.
# These variants give extra information to the type checker.
# They are ignored at runtime.
ClickEvent, DragEvent = Any


@overload
def mouse_event(x1: int, y1: int) -> ClickEvent:
    ...


@overload
def mouse_event(x1: int, y1: int, x2: int, y2: int) -> DragEvent:
    ...


# The actual *implementation* of 'mouse_event'.
# The implementation contains the actual runtime logic.
#
# It may or may not have type hints. If it does, mypy
# will check the body of the implementation against the
# type hints.
#
# Mypy will also check and make sure the signature is
# consistent with the provided variants.


def mouse_event(
    x1: int, y1: int, x2: Optional[int] = None, y2: Optional[int] = None
) -> Union[ClickEvent, DragEvent]:
    if x2 is None and y2 is None:
        return ClickEvent(x1, y1)
    elif x2 is not None and y2 is not None:
        return DragEvent(x1, y1, x2, y2)
    else:
        raise TypeError("Bad arguments")
