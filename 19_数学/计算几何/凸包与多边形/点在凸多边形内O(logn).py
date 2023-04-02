from typing import List, Tuple


def cross(a: Tuple[int, int], b: Tuple[int, int], c: Tuple[int, int]) -> int:
    return (b[0] - a[0]) * (c[1] - a[1]) - (b[1] - a[1]) * (c[0] - a[0])


def polygonAndPoints(poly: List[Tuple[int, int]], p: Tuple[int, int]) -> int:
    """O(logn) 判断点p与凸多边形poly的关系.其中多边形的顶点按逆时针方向给出.

    Returns:
      -1: 在多边形外部
      0: 在多边形边上
      1: 在多边形内部

    Notes:
      1. 二分在哪个三角形内部;
      2. 然后判断点是否在线段上, 最后内积定位在里面还是外面.
    """
    n = len(poly)
    left = 1
    right = n
    p0 = poly[0]
    while left + 1 < right:
        mid = (left + right) >> 1
        if cross(p0, p, poly[mid]) <= 0:
            left = mid
        else:
            right = mid
    if left == n - 1:
        left -= 1

    pi = poly[left]
    pj = poly[left + 1]
    c0 = cross(p0, pi, pj)
    c1 = cross(p0, p, pj)
    c2 = cross(p0, pi, p)
    if c0 < 0:
        c1 = -c1
        c2 = -c2
    if 0 <= c1 and 0 <= c2 and c1 + c2 <= c0:
        cur = cross(p, pi, pj)
        if pi == poly[1]:
            cur *= cross(p, p0, pi)
        if pj == poly[-1]:
            cur *= cross(p, p0, pj)
        return 0 if cur == 0 else 1
    return -1


if __name__ == "__main__":
    # https://atcoder.jp/contests/abc296/tasks/abc296_g
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    poly = list(tuple(map(int, input().split())) for _ in range(n))
    q = int(input())
    for _ in range(q):
        x, y = map(int, input().split())
        res = polygonAndPoints(poly, (x, y))
        if res == -1:
            print("OUT")
        elif res == 0:
            print("ON")
        else:
            print("IN")
