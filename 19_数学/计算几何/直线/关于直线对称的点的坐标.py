# 点の線対称
# 关于直线对称的点的坐标
# https://tjkendev.github.io/procon-library/python/geometry/reflection_point.html
from typing import Tuple

# 这里的直线是两点式


def reflection(
    line: Tuple[Tuple[int, int], Tuple[int, int]], point: Tuple[int, int]
) -> Tuple[float, float]:
    def dot3(a, b, c):
        x0, y0 = a
        x1, y1 = b
        x2, y2 = c
        return (x1 - x0) * (x2 - x0) + (y1 - y0) * (y2 - y0)

    def dist2(a, b):
        x0, y0 = a
        x1, y1 = b
        return (x0 - x1) * (x0 - x1) + (y0 - y1) * (y0 - y1)

    p0, p1 = line
    dv = dot3(p0, p1, point)
    dd = dist2(p0, p1)

    x0, y0 = p0
    x1, y1 = p1
    xp, yp = point

    xm = x0 + dv * (x1 - x0) / dd
    ym = y0 + dv * (y1 - y0) / dd

    return xp + 2 * (xm - xp), yp + 2 * (ym - yp)


# # example
# print(reflection(((0, 2), (2, 1)), (0, 0)))
# # => (1.6, 3.2)

if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    x1, y1, x2, y2 = map(int, input().split())
    line = ((x1, y1), (x2, y2))
    q = int(input())
    for _ in range(q):
        x, y = map(int, input().split())
        point = (x, y)
        print(*reflection(line, point))
