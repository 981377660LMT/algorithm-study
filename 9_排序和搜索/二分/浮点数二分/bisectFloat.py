from math import ceil
from typing import Callable


def bisectLeftFloat(
    left: float, right: float, check: Callable[[float], bool], absErrorInv=int(1e9)
) -> float:
    diff = ceil((right - left) * absErrorInv)
    round = diff.bit_length()
    for _ in range(round):
        mid = (left + right) / 2
        if check(mid):
            right = mid
        else:
            left = mid
    return (left + right) / 2


def bisectRightFloat(
    left: float, right: float, check: Callable[[float], bool], absErrorInv=int(1e9)
) -> float:
    diff = ceil((right - left) * absErrorInv)
    round = diff.bit_length()
    for _ in range(round):
        mid = (left + right) / 2
        if check(mid):
            left = mid
        else:
            right = mid
    return (left + right) / 2
