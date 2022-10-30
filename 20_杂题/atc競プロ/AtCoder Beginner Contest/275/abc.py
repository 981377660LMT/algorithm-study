import sys

# 二次元平面があります。1 以上 9 以下の整数 r,c について、S
# r
# ​
#   の c 番目の文字が # であるとき座標 (r,c) にポーンが置いてあり、S
# r
# ​
#   の c 番目の文字が . であるとき座標 (r,c) に何も置かれていません。

# この平面上の正方形であって、4 頂点全てにポーンが置いてあるものの個数を求めてください。

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

from itertools import combinations
from typing import List


def dist(a: List[int], b: List[int]) -> int:
    return (a[0] - b[0]) ** 2 + (a[1] - b[1]) ** 2


class Solution:
    def validSquare(self, p1: List[int], p2: List[int], p3: List[int], p4: List[int]) -> bool:
        """边长不为0 四条边相等 对角线相等"""
        dists = sorted([dist(a, b) for a, b in combinations((p1, p2, p3, p4), 2)])
        d1, d2, d3, d4, d5, d6 = dists
        return d1 > 0 and d1 == d2 == d3 == d4 and d5 == d6


if __name__ == "__main__":
    S = []
    for r in range(9):
        row = input()
        for c in range(9):
            if row[c] == "#":
                S.append((r, c))
    res = 0
    for p1, p2, p3, p4 in combinations(S, 4):
        if Solution().validSquare(p1, p2, p3, p4):
            res += 1
    print(res)
