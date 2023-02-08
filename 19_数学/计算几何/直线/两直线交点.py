# 两直线交点
# https://tjkendev.github.io/procon-library/python/geometry/line_cross_point.html

from typing import Tuple, Union

Point = Tuple[int, int]


def line_cross_point1(
    p0: Point, p1: Point, q0: Point, q1: Point
) -> Union[Tuple[float, float], Tuple[None, None]]:
    """两点式直线求交点"""
    x0, y0 = p0
    x1, y1 = p1
    x2, y2 = q0
    x3, y3 = q1
    a0 = x1 - x0
    b0 = y1 - y0
    a2 = x3 - x2
    b2 = y3 - y2

    d = a0 * b2 - a2 * b0
    if d == 0:
        # two lines are parallel
        return None, None

    sn = b2 * (x2 - x0) - a2 * (y2 - y0)
    return x0 + a0 * sn / d, y0 + b0 * sn / d


def line_cross_point2(
    a1: float, b1: float, c1: float, a2: float, b2: float, c2: float
) -> Union[Tuple[float, float], Tuple[None, None]]:
    """直线表达式为ax+by=c  求两直线交点"""
    if a1 * b2 == a2 * b1:  # 平行
        return None, None
    x = (c1 * b2 - c2 * b1) / (a1 * b2 - a2 * b1)
    y = (c2 * a1 - c1 * a2) / (a1 * b2 - a2 * b1)
    return x, y


if __name__ == "__main__":
    q = int(input())
    for _ in range(q):
        x0, y0, x1, y1, x2, y2, x3, y3 = map(int, input().split())
        cx, cy = line_cross_point1((x0, y0), (x1, y1), (x2, y2), (x3, y3))
        print("%.16f %.16f" % (cx, cy))
