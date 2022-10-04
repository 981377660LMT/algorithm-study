"""中垂线/垂直平分线"""
# 等腰梯形 Isosceles Trapezium
# !求选出的等腰梯形的顶点和的最大值
# n<=1000

# O(n^2logn)
# !等腰梯形的平行边:
# !1.中垂线相等
# !2.中点不等
from collections import defaultdict
from functools import reduce
from heapq import nlargest
from itertools import combinations
from math import gcd
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def center2(x1: int, y1: int, x2: int, y2: int):
    return (x1 + x2), (y1 + y2)


def 垂直平分线(x1: int, y1: int, x2: int, y2: int):
    """求垂直平分线的方程ax+by+c=0"""
    a = (x2 - x1) * 2
    b = (y2 - y1) * 2
    c = x1 * x1 - x2 * x2 + y1 * y1 - y2 * y2
    g = reduce(gcd, [a, b, c])
    if a < 0 or (a == 0 and b < 0):
        g *= -1
    return a // g, b // g, c // g


if __name__ == "__main__":
    n = int(input())
    points = []
    for _ in range(n):
        x, y, weight = map(int, input().split())
        points.append((x, y, weight))

    counter = defaultdict(lambda: defaultdict(int))
    for p1, p2 in combinations(points, 2):
        x1, y1, w1 = p1
        x2, y2, w2 = p2
        line = 垂直平分线(x1, y1, x2, y2)
        center = center2(x1, y1, x2, y2)
        counter[line][center] = max(counter[line][center], w1 + w2)

    res = -1
    for v in counter.values():
        if len(v) >= 2:
            max2 = nlargest(2, v.values())
            res = max(res, sum(max2))
    print(res)
