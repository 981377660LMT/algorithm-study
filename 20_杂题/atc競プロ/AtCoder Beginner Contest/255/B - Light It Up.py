# !求每个人至少被一个灯照亮的最小光强
# !对每个人寻找最近的灯泡

from math import sqrt
import sys
import os

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def main() -> None:
    n, k = map(int, input().split())
    lights = [int(num) - 1 for num in input().split()]
    points = []
    for _ in range(n):
        r, c = map(int, input().split())
        points.append((r, c))

    res = -int(1e18)
    for i in range(n):
        minDist = int(1e18)
        r1, c1 = points[i]
        for j in lights:
            r2, c2 = points[j]
            minDist = min(minDist, sqrt((r1 - r2) ** 2 + (c1 - c2) ** 2))
        res = max(res, minDist)
    print(res)


if os.environ.get("USERNAME", "") == "caomeinaixi":
    while True:
        try:
            main()
        except (EOFError, ValueError):
            break
else:
    main()
