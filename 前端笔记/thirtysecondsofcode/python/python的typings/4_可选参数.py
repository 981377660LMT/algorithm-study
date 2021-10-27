from typing import Sequence, Optional


def player_order(names: Sequence[str], start: Optional[str] = None) -> Sequence[str]:
    ...


# 等价于Union类型的 Union[None, str]，意思是这个参数的值类型为str，默认的话可以是None
# 在 3.10 版更改: Optional can now be written as X | None. See union type expressions.


def concat(x: Optional[str], y: Optional[str]) -> Optional[str]:
    if x is not None and y is not None:
        # Both x and y are not None here
        return x + y
    else:
        return None

