# 平面最远点对-凸多边形的直径
# Rotating Calipers
# https://tjkendev.github.io/procon-library/python/geometry/rotating_calipers.html
# キャリパー法(Rotating Calipers法)
# 凸多角形を回転しながら走査し、最遠点対を求める

# Diameter of a Convex Polygon


from math import sqrt
from typing import List, Tuple


def diameterConvexPolygon(ps: List[Tuple[int, int]]) -> float:
    """求出凸多边形的直径(最远点对),其中多边形的顶点按逆时针方向给出."""
    ps = sorted(ps)
    qs = convex_hull(ps)
    n = len(qs)
    if n == 2:
        return dist(qs[0], qs[1])
    i = j = 0
    for k in range(n):
        if qs[k] < qs[i]:
            i = k
        if qs[j] < qs[k]:
            j = k
    res = 0
    si = i
    sj = j
    while i != sj or j != si:
        cand = dist(qs[i], qs[j])
        if cand > res:
            res = cand
        if cross(qs[i], qs[i - n + 1], qs[j], qs[j - n + 1]) < 0:
            i = (i + 1) % n
        else:
            j = (j + 1) % n
    return res


def cross(a, b, c, d):
    return (b[0] - a[0]) * (d[1] - c[1]) - (b[1] - a[1]) * (d[0] - c[0])


def convex_hull(ps):
    qs = []
    n = len(ps)
    for p in ps:
        while len(qs) > 1 and cross(qs[-1], qs[-2], qs[-1], p) >= 0:
            qs.pop()
        qs.append(p)
    t = len(qs)
    for i in range(n - 2, -1, -1):
        p = ps[i]
        while len(qs) > t and cross(qs[-1], qs[-2], qs[-1], p) >= 0:
            qs.pop()
        qs.append(p)
    return qs


def dist(a, b):
    return sqrt((a[0] - b[0]) * (a[0] - b[0]) + (a[1] - b[1]) * (a[1] - b[1]))


if __name__ == "__main__":
    n = int(input())
    points = list(tuple(map(float, input().split())) for _ in range(n))
    print("%.09f" % diameterConvexPolygon(points))
