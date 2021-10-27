#  Python 运行时不强制执行函数和变量类型注解，但这些注解可用于类型检查器、IDE、静态检查器等第三方工具。
#  类型别名（type alias）
from typing import Dict, List, Tuple
from collections.abc import Sequence


Vector = List[float]


def scale(scalar: float, vector: Vector) -> Vector:
    return [scalar * num for num in vector]


# typechecks; a list of floats qualifies as a Vector.
new_vector = scale(2.0, [1.0, -4.2, 5.4])

############################################################

ConnectionOptions = Dict[str, str]
Address = Tuple[str, int]
Server = Tuple[Address, ConnectionOptions]


def broadcast_message(message: str, servers: Sequence[Server]) -> None:
    ...


# The static type checker will treat the previous type signature as
# being exactly equivalent to this one.
def broadcast_message2(
    message: str, servers: Sequence[Tuple[Tuple[str, int], Dict[str, str]]]
) -> None:
    ...
