from fractions import Fraction
from typing import List


Point = List[int]


def calCross(p1: Point, p2: Point, p3: Point) -> int:
    ab = [p2[0] - p1[0], p2[1] - p1[1]]
    ac = [p3[0] - p1[0], p3[1] - p1[1]]
    return ab[0] * ac[1] - ab[1] * ac[0]


def calTriangleArea(p1: Point, p2: Point, p3: Point) -> Fraction:
    """三角形面积公式"""
    return Fraction(abs(calCross(p1, p2, p3)), 2)


def calTriangleArea2(p1: Point, p2: Point, p3: Point) -> int:
    """三角形面积公式2倍"""
    return abs(calCross(p1, p2, p3))


def calPolygenArea(points: List[Point]) -> Fraction:
    """多边形面积公式 鞋带定理"""
    assert len(points) >= 3
    res = Fraction(0)
    for (x1, y1), (x2, y2) in zip(points, points[1:] + [points[0]]):
        res += x1 * y2 - x2 * y1
    return Fraction(abs(res), 2)


def calPolygenArea2(points: List[Point]) -> int:
    """多边形面积的2倍 鞋带定理"""
    assert len(points) >= 3
    res = 0
    for (x1, y1), (x2, y2) in zip(points, points[1:] + [points[0]]):
        res += x1 * y2 - x2 * y1
    return abs(res)
