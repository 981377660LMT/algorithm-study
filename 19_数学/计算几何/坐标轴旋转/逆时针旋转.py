from math import cos, radians, sin
from typing import Tuple


def rotate(x: int, y: int, *, deg: int) -> Tuple[float, float]:
    """(x, y)逆时针旋转deg度(角度制)后的坐标"""
    rad = radians(deg)
    return x * cos(rad) - y * sin(rad), x * sin(rad) + y * cos(rad)
