"""线性规划 求可行域 可行解"""

# a1x+b1y<=c1
# a2x+b2y<=c2
# ...
# anx+bny<=cn

# n<=500
# ai,bi,ci<=1e9

# !求 (x+y) 最大值
# !直线两两的交点处的最大值
# O(n^3)求解
from itertools import combinations
from typing import Tuple
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


def solve(
    a1: float, b1: float, c1: float, a2: float, b2: float, c2: float
) -> Tuple[float, float] | Tuple[None, None]:
    """直线表达式为ax+by=c  求两直线交点"""
    if a1 * b2 == a2 * b1:  # 平行
        return None, None
    x = (c1 * b2 - c2 * b1) / (a1 * b2 - a2 * b1)
    y = (c2 * a1 - c1 * a2) / (a1 * b2 - a2 * b1)
    return x, y


def check(x: float, y: float) -> bool:
    return all(a * x + b * y <= c for a, b, c in lines)


n = int(input())
lines = []
for _ in range(n):
    a, b, c = map(int, input().split())
    lines.append([a, b, c])
res = -int(1e20)

for line1, line2 in combinations(lines, 2):
    x, y = solve(*line1, *line2)
    if x is None or y is None:
        continue
    if check(x, y):
        res = max(res, x + y)
print(res)
