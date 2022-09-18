"""计算圆心和半径 (x,y,r)"""

from typing import Tuple, Union


def calCircle1(
    x1: int, y1: int, x2: int, y2: int, r: int
) -> Union[Tuple[None, None, None], Tuple[float, float, float]]:
    """已知两点坐标(x1,y1),(x2,y2)和半径r,求圆的圆心坐标(x,y)和半径r

    可以确定两个圆 这里返回的是`位于右半平面`的圆心
    """
    d = ((x1 - x2) * (x1 - x2) + (y1 - y2) * (y1 - y2)) / 4.0
    if d > r * r:
        return None, None, None
    x0 = (x1 + x2) / 2.0 + (y2 - y1) * (r * r - d) ** 0.5 / (d * 4) ** 0.5
    y0 = (y1 + y2) / 2.0 - (x2 - x1) * (r * r - d) ** 0.5 / (d * 4) ** 0.5
    return x0, y0, r


def calCircle2(
    x1: int, y1: int, x2: int, y2: int, x3: int, y3: int
) -> Union[Tuple[None, None, None], Tuple[float, float, float]]:
    """三点圆公式,求圆的圆心坐标(x,y)和半径r"""
    a, b, c, d = x1 - x2, y1 - y2, x1 - x3, y1 - y3
    a1 = (x1 * x1 - x2 * x2 + y1 * y1 - y2 * y2) / 2
    a2 = (x1 * x1 - x3 * x3 + y1 * y1 - y3 * y3) / 2
    theta = b * c - a * d
    if theta == 0:
        return None, None, None
    x0 = (b * a2 - d * a1) / theta
    y0 = (c * a1 - a * a2) / theta
    return x0, y0, ((x1 - x0) * (x1 - x0) + (y1 - y0) * (y1 - y0)) ** 0.5


if __name__ == "__main__":
    x, y, r = calCircle1(0, 0, 1, 1, 1)
    print(x, y, r)
    x, y, r = calCircle2(0, 0, 2, 0, 1, 2)
    print(x, y, r)
