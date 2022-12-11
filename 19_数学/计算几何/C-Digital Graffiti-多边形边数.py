"""
求黑色多边形边数

!边数=角点数,只需统计四个相邻点,三白一黑或者三黑一白即为边界角点
"""

import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# '#' 代表黑色，'.' 代表白色。
# 给定一个包含黑色多边形的矩阵,求黑色多边形的边数。


def countEdge(grid: List[str]) -> int:
    ROW, COL = len(grid), len(grid[0])
    res = 0
    for r in range(ROW - 1):
        for c in range(COL - 1):
            neighbors = [grid[r][c], grid[r][c + 1], grid[r + 1][c], grid[r + 1][c + 1]]
            if neighbors.count("#") == 3 or neighbors.count("#") == 1:
                res += 1
    return res


if __name__ == "__main__":
    ROW, COL = map(int, input().split())
    grid = [input() for _ in range(ROW)]
    print(countEdge(grid))
