"""给定二维空间中四点的坐标，返回四点是否可以构造一个矩形。"""
# 验证矩形/矩形判定
# 判断矩形


from typing import Tuple


Point = Tuple[int, int]


def isRectangleAnyOrder(p1: Point, p2: Point, p3: Point, p4: Point) -> bool:
    """判断四个点是否可以构成一个矩形"""
    return (
        _isRectangle(p1, p2, p3, p4) or _isRectangle(p1, p2, p4, p3) or _isRectangle(p1, p3, p2, p4)
    )


def _isRectangle(p1: Point, p2: Point, p3: Point, p4: Point) -> bool:
    return _isOrthogonal(p1, p2, p3) and _isOrthogonal(p2, p3, p4) and _isOrthogonal(p3, p4, p1)


def _isOrthogonal(a: Point, b: Point, c: Point) -> bool:
    """判断 ∠abc 是否为直角"""
    newP1 = (a[0] - b[0], a[1] - b[1])
    newP2 = (c[0] - b[0], c[1] - b[1])
    return newP1[0] * newP2[0] + newP1[1] * newP2[1] == 0  # dot(newP1, newP2) == 0
