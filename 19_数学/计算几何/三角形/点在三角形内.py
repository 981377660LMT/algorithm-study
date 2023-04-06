# 点在三角形内

from typing import Tuple


def inTriangle(
    a: Tuple[int, int], b: Tuple[int, int], c: Tuple[int, int], p: Tuple[int, int]
) -> bool:
    """判断点是否在三角形内"""

    pa, pb, pc = (a[0] - p[0], a[1] - p[1]), (b[0] - p[0], b[1] - p[1]), (c[0] - p[0], c[1] - p[1])
    res1 = abs(det((b[0] - a[0], b[1] - a[1]), (c[0] - a[0], c[1] - a[1])))
    res2 = abs(det(pa, pb)) + abs(det(pb, pc)) + abs(det(pc, pa))
    return res1 == res2


def det(a: Tuple[int, int], b: Tuple[int, int]) -> int:
    return a[0] * b[1] - a[1] * b[0]
