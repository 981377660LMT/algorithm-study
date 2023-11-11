# 求正多边形(正n边形)的顶点坐标 其中n为偶数

# 给定的两个点连线的中点就是正n 边形外切圆的圆心，
# 所以我们可以直接算出圆心，然后旋转2pi/n 度

from math import cos, pi, sin
from typing import Tuple


def opposite(n: int, x0: int, y0: int, xmid: int, ymid: int) -> Tuple[float, float]:
    centerX, centerY = (x0 + xmid) / 2, (y0 + ymid) / 2
    dx, dy = x0 - centerX, y0 - centerY
    theta = 2 * pi / n
    return centerX + dx * cos(theta) - dy * sin(theta), centerY + dx * sin(theta) + dy * cos(theta)


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n = int(input())
    x0, y0 = map(int, input().split())
    xmid, ymid = map(int, input().split())
    x1, y1 = opposite(n, x0, y0, xmid, ymid)
    print(x1, y1)
