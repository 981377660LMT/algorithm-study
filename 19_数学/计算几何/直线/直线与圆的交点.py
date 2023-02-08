# 直线与圆的交点

from math import sqrt
from typing import List, Tuple

# !直线与圆的交点
# (cx, cy, r) : 圆心坐标和半径
# (x1, y1, x2, y2) : 直线上的两点
def crossPointsOfCiclreAndPoint(
    cx: int, cy: int, r: int, x1: int, y1: int, x2: int, y2: int
) -> List[Tuple[float, float]]:
    xd = x2 - x1
    yd = y2 - y1
    X = x1 - cx
    Y = y1 - cy
    a = xd * xd + yd * yd
    b = xd * X + yd * Y
    c = X * X + Y * Y - r * r
    D = b * b - a * c  # D = 0の時は1本で、D < 0の時は存在しない
    if D < 0:
        return []
    s1 = (-b + sqrt(D)) / a
    s2 = (-b - sqrt(D)) / a
    if s1 == s2:
        return [(x1 + xd * s1, y1 + yd * s1)]
    return [(x1 + xd * s1, y1 + yd * s1), (x1 + xd * s2, y1 + yd * s2)]


if __name__ == "__main__":
    # https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_D&lang=ja
    cx, cy, r = map(int, input().split())
    q = int(input())
    for _ in range(q):
        x1, y1, x2, y2 = map(int, input().split())
        res = sorted(crossPointsOfCiclreAndPoint(cx, cy, r, x1, y1, x2, y2))
        if len(res) == 1:
            p0, p1 = res[0], res[0]
        else:
            p0, p1 = res[0], res[1]
        print("%.08f %.08f %.08f %.08f" % (p0 + p1))
