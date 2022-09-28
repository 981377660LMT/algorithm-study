# 平面上n点组成多少个三角形/三角形的个数 O(n^2)
# !comb(n,3)再减去三点共线的对数


from collections import defaultdict
from math import comb, gcd
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


def calSlope1(x1: int, y1: int, x2: int, y2: int) -> Tuple[int, int]:
    """直线斜率"""
    if x2 == x1:
        return (INF, INF)
    gcd_ = gcd(x2 - x1, y2 - y1)
    a, b = (y2 - y1) // gcd_, (x2 - x1) // gcd_
    if a == 0:
        return (0, b if b > 0 else -b)
    elif a < 0:
        return (-a, -b)
    else:
        return (a, b)


def solve(points: List[Tuple[int, int]]) -> int:
    """平面上n点组成多少个三角形 O(n^2)"""
    n = len(points)
    res = 0  # 三点共线的对数
    for i in range(n):
        x1, y1 = points[i]
        slopeCounter = defaultdict(int)
        for j in range(i + 1, n):
            x2, y2 = points[j]
            slope = calSlope1(x1, y1, x2, y2)
            slopeCounter[slope] += 1

        for count in slopeCounter.values():
            res += comb(count, 2)

    return comb(n, 3) - res


if __name__ == "__main__":
    n = int(input())
    points = [tuple(map(int, input().split())) for _ in range(n)]
    print(solve(points))
