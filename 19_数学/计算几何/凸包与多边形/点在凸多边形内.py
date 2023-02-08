from typing import List, Tuple

Point = Tuple[int, int]


def cross(a, b, c):
    return (b[0] - a[0]) * (c[1] - a[1]) - (b[1] - a[1]) * (c[0] - a[0])


def inside_convex_polygon(p0: Point, poly: List[Point]) -> bool:
    """O(log N) 判断点是否在凸多边形内(包括边界),其中多边形的顶点按逆时针方向给出."""
    L = len(poly)
    left = 1
    right = L
    q0 = poly[0]
    while left + 1 < right:
        mid = (left + right) >> 1
        if cross(q0, p0, poly[mid]) <= 0:
            left = mid
        else:
            right = mid
    if left == L - 1:
        left -= 1
    qi = poly[left]
    qj = poly[left + 1]
    v0 = cross(q0, qi, qj)
    v1 = cross(q0, p0, qj)
    v2 = cross(q0, qi, p0)
    if v0 < 0:
        v1 = -v1
        v2 = -v2
    return 0 <= v1 and 0 <= v2 and v1 + v2 <= v0


# O(N)
def _inside_convex_polygon2(p0: Point, poly: List[Point]) -> bool:
    L = len(poly)
    D = [cross(poly[i - 1], p0, poly[i]) for i in range(L)]
    return all(e >= 0 for e in D) or all(e <= 0 for e in D)


if __name__ == "__main__":
    qs = [(-2, 0), (0, -2), (2, 0), (0, 2)]
    print(inside_convex_polygon((0, 0), qs))
    # => "True"
    print(inside_convex_polygon((1, 1), qs))
    # => "True"
    print(inside_convex_polygon((2, 2), qs))
    # => "False"
