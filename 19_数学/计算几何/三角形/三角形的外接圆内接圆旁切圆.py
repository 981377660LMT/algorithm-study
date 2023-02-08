# 三角形的外接圆/内接圆/旁切圆

from typing import List, Tuple


Point = Tuple[int, int]


def circumcircle(P1: "Point", P2: "Point", P3: "Point") -> Tuple[float, float, float]:
    """三角形的外接圆"""
    x1, y1 = P1
    x2, y2 = P2
    x3, y3 = P3
    a = 2 * (x1 - x2)
    b = 2 * (y1 - y2)
    p = x1**2 - x2**2 + y1**2 - y2**2
    c = 2 * (x1 - x3)
    d = 2 * (y1 - y3)
    q = x1**2 - x3**2 + y1**2 - y3**2
    det = a * d - b * c
    x = d * p - b * q
    y = a * q - c * p
    if det < 0:
        x = -x
        y = -y
        det = -det
    x /= det
    y /= det
    r = ((x - x1) ** 2 + (y - y1) ** 2) ** 0.5
    return x, y, r


def incircle(P1: "Point", P2: "Point", P3: "Point") -> Tuple[float, float, float]:
    """三角形的内接圆"""
    x1, y1 = P1
    x2, y2 = P2
    x3, y3 = P3

    dx1 = x2 - x1
    dy1 = y2 - y1
    dx2 = x3 - x1
    dy2 = y3 - y1

    d1 = ((x3 - x2) ** 2 + (y3 - y2) ** 2) ** 0.5
    d2 = (dx2**2 + dy2**2) ** 0.5
    d3 = (dx1**2 + dy1**2) ** 0.5
    dsum = d1 + d2 + d3

    r = abs(dx1 * dy2 - dx2 * dy1) / dsum
    x = (x1 * d1 + x2 * d2 + x3 * d3) / dsum
    y = (y1 * d1 + y2 * d2 + y3 * d3) / dsum
    return x, y, r


def excircle(P1: "Point", P2: "Point", P3: "Point") -> List[Tuple[float, float, float]]:
    """三角形的三个旁切圆

    分别是与 p2-p3边相切的圆, p1-p3边相切的圆, p1-p2边相切的圆
    """
    x1, y1 = P1
    x2, y2 = P2
    x3, y3 = P3

    dx1 = x2 - x1
    dy1 = y2 - y1
    dx2 = x3 - x1
    dy2 = y3 - y1

    d1 = ((x3 - x2) ** 2 + (y3 - y2) ** 2) ** 0.5
    d3 = (dx1**2 + dy1**2) ** 0.5
    d2 = (dx2**2 + dy2**2) ** 0.5

    S2 = abs(dx1 * dy2 - dx2 * dy1)

    dsum1 = -d1 + d2 + d3
    r1 = S2 / dsum1
    ex1 = (-x1 * d1 + x2 * d2 + x3 * d3) / dsum1
    ey1 = (-y1 * d1 + y2 * d2 + y3 * d3) / dsum1

    dsum2 = d1 - d2 + d3
    r2 = S2 / dsum2
    ex2 = (x1 * d1 - x2 * d2 + x3 * d3) / dsum2
    ey2 = (y1 * d1 - y2 * d2 + y3 * d3) / dsum2

    dsum3 = d1 + d2 - d3
    r3 = S2 / dsum3
    ex3 = (x1 * d1 + x2 * d2 - x3 * d3) / dsum3
    ey3 = (y1 * d1 + y2 * d2 - y3 * d3) / dsum3

    return [(ex1, ey1, r1), (ex2, ey2, r2), (ex3, ey3, r3)]
