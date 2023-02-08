from typing import List, Tuple


Point = Tuple[int, int]


def line_polygon_intersection(p0: Point, p1: Point, qs: List[Point]) -> List[int]:
    """
    O(log N) 求直线与凸多边形的交点在多边形上的位置ei(哪条边上)
    ei 表示vi-1, vi 之间的边

    其中多边形的顶点按逆时针方向给出.
    """
    x0, y0 = p0
    x1, y1 = p1
    dx = x1 - x0
    dy = y1 - y0
    h = lambda p: (p[0] - x0) * dy - (p[1] - y0) * dx
    L = len(qs)

    i0 = i1 = -1

    v0 = h(qs[0])
    v1 = h(qs[L - 1])
    if v0 == v1:
        v2 = h(qs[1])
        # assert v0 != v2
        if v0 < v2:
            i0 = L - 1
        else:
            i1 = L - 1
    else:
        v2 = h(qs[1])
        if v1 > v0 <= v2:
            i0 = 0
        elif v1 < v0 >= v2:
            i1 = 0
        else:
            g = lambda x: min((v1 - v0) * x / (L - 1) + v0, h(qs[x]))
            i0 = binary_search(lambda x: g(x + 1) - g(x), L - 1)

    if i1 == -1:
        B = i0 - L
        k = binary_search(lambda x: h(qs[B + x]) - h(qs[B + x + 1]), L)
        i1 = (i0 + k) % L
    else:
        B = i1 - L
        k = binary_search(lambda x: h(qs[B + x + 1]) - h(qs[B + x]), L)
        i0 = (i1 + k) % L

    if h(qs[i0]) * h(qs[i1]) > 0:
        # a line and a polygon are disjoint
        return []

    # a vertex to the left side of a line: i0
    # a vertex to the right side of a line: i1

    f = lambda i: h(qs[i - L])
    k0 = find_zero(i1, i0 if i1 < i0 else i0 + L, f) % L
    k1 = find_zero(i0, i1 if i0 < i1 else i1 + L, f) % L
    # vertices to the left side of a line: k0, k0+1, ..., k1-2, k1-1
    # vertices to the right side of a line: k1, k1+1, ..., k0-2, k0-1
    if k0 == k1:
        return [k0]
    return [k0, k1]


# 直线与凸多边形的交点 O(logn)
def find_zero(x0, x1, f):
    v0 = f(x0)
    if v0 == 0:
        return x0 + 1
    left = x0
    right = x1 + 1
    while left + 1 < right:
        mid = (left + right) >> 1
        if v0 * f(mid) >= 0:
            left = mid
        else:
            right = mid
    return right


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


if __name__ == "__main__":
    # e_i is a line segment between v_{i-1} and v_i
    p0, p1 = (0, 0), (1, 1)
    qs = [(0, 2), (2, 0), (4, 0), (4, 2), (2, 4), (0, 4)]
    print(line_polygon_intersection(p0, p1, qs))
    # => "[4, 1]": cross points is on either e_4 or e_1

    p0, p1 = (0, -1), (4, 0)
    qs = [(0, 2), (2, 0), (4, 0), (4, 2), (2, 4), (0, 4)]
    # => "[3]": cross point is on e_3
