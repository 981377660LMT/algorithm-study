"""
平面上有n个不同的点
求通过>=k个点的直线的条数
n,k<=300
存在无数个输出'Infinity'
"""

from functools import lru_cache
from typing import Tuple, Union
from collections import defaultdict
from fractions import Fraction
from math import gcd
import sys
import os

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
gcd = lru_cache(maxsize=None)(gcd)


def calSlopeInterceptForm(
    x1: int, y1: int, x2: int, y2: int
) -> Union[Tuple[Fraction, Fraction], Tuple[None, int]]:
    """求出直线方程的斜截式`y=kx+b 或 x=b`的斜率和截距

    Returns:
        Tuple[Fraction, Fraction] | Tuple[None, int]: 斜率和截距
        当斜率为None时,截距表示在x轴上的截距,直线方程为 x=截距;
        当斜率不为None时,截距表示在y轴上的截距,直线方程为 y=斜率*x+截距
    """
    if x1 == x2:
        return None, x1
    slope = Fraction(y2 - y1, x2 - x1)
    intercept = y1 - slope * x1
    return slope, intercept


def calGeneralForm(x1: int, y1: int, x2: int, y2: int) -> Tuple[int, int, int]:
    """求出直线方程的一般式`ax+by+c=0`的系数"""
    if x1 == x2:
        return 1, 0, -x1
    dy, dx = y2 - y1, x2 - x1
    gcd_ = gcd(dy, dx)
    dy, dx = dy // gcd_, dx // gcd_
    c = dx * y1 - dy * x1
    if dx > 0:
        dx, dy, c = -dx, -dy, -c
    return dy, -dx, c


def main() -> None:
    n, k = map(int, input().split())
    if k == 1:
        print("Infinity")
        exit(0)
    points = []
    for _ in range(n):
        x, y = map(int, input().split())
        points.append((x, y))

    res = set()  # 直线方程的集合
    for i in range(n):
        x1, y1 = points[i]
        slopeGroup = defaultdict(list)
        for j in range(i + 1, n):
            x2, y2 = points[j]
            if x1 == x2:
                slopeGroup[None].append(j)
            else:
                slope = Fraction(y2 - y1, x2 - x1)
                slopeGroup[slope].append(j)
        for slope, group in slopeGroup.items():
            if len(group) + 1 >= k:
                x2, y2 = points[group[0]]
                res.add(calGeneralForm(x1, y1, x2, y2))
    print(len(res))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
