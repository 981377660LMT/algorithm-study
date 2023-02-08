# tangent lines from a point to a convex polygon
# 求多边形外一点出发到该多边形的切线的交点
# 就是说一条直线与 多边形 相交，并且 多边形 在直线的一侧，这样的直线称之为多边形的切线
# O(logn)

from math import atan2
from typing import Tuple

EPS = 1e-12


def convex_polygon_tangent(p0, qs) -> Tuple[int, int]:
    L = len(qs)
    d = L // 3
    gx = (qs[0][0] + qs[d][0] + qs[2 * d][0]) / 3
    gy = (qs[0][1] + qs[d][1] + qs[2 * d][1]) / 3

    x0, y0 = p0
    dx = gx - x0
    dy = gy - y0
    sv = -x0 * dy + y0 * dx
    cv = -x0 * dx - y0 * dy
    h = lambda p: atan2(p[0] * dy - p[1] * dx + sv, p[0] * dx + p[1] * dy + cv)

    i0 = i1 = -1

    v0 = h(qs[0])
    v1 = h(qs[L - 1])
    if abs(v0 - v1) < EPS:
        v2 = h(qs[1])
        # assert v2 != v0
        if v0 < v2:
            i0 = L - 1
        else:
            i1 = L - 1
    else:
        v2 = h(qs[1])
        if v2 >= v0 < v1:
            i0 = 0
        elif v2 <= v0 > v1:
            i1 = 0
        else:
            g = lambda x: min((v1 - v0) * x / (L - 1) + v0, h(qs[x]))
            i0 = binary_search(lambda x: g(x + 1) - g(x), L - 1)

    if i1 == -1:
        B = i0 - L
        k = binary_search(lambda x: h(qs[x + B]) - h(qs[x + B + 1]), L)
        i1 = (i0 + k) % L
    else:
        B = i1 - L
        k = binary_search(lambda x: h(qs[x + B + 1]) - h(qs[x + B]), L)
        i0 = (i1 + k) % L
    # assert i0 != -1 != i1
    return i0, i1


def binary_search(f, L):
    left = 0
    right = L
    while left + 1 < right:
        mid = (left + right) >> 1
        if f(mid) < 0:
            left = mid
        else:
            right = mid
    return right
