from typing import List


Point = List[int]


def calCross(p1: Point, p2: Point, p3: Point) -> int:
    ab = [p2[0] - p1[0], p2[1] - p1[1]]
    ac = [p3[0] - p1[0], p3[1] - p1[1]]
    return ab[0] * ac[1] - ab[1] * ac[0]


def calTriangleArea(p1: Point, p2: Point, p3: Point):
    """三角形面积公式"""
    return abs(calCross(p1, p2, p3)) / 2


def calArea(points: List[Point]) -> float:
    """凸多边形面积公式
    
    以多边形的某一点为顶点,将其划分成几个三角形,计算这些三角形的面积,然后加起来
    """
    assert len(points) >= 3
    res = 0
    for i in range(1, len(points) - 1):
        res += calTriangleArea(points[0], points[i], points[i + 1])
    return res
