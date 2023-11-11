# 求出黑色区域的范围 返回row1,row2,col1,col2
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(grid: List[List[int]]) -> List[int]:
    """求出全1矩形边界(四个顶点) 返回row1,row2,col1,col2"""
    ROW, COL = len(grid), len(grid[0])
    row1, row2, col1, col2 = INF, -INF, INF, -INF
    for r in range(ROW):
        for c in range(COL):
            if grid[r][c] == 1:
                row1, row2, col1, col2 = min(row1, r), max(row2, r), min(col1, c), max(col2, c)
    return [row1, row2, col1, col2]


if __name__ == "__main__":
    grid = []
    for _ in range(10):
        row = input()
        grid.append([1 if c == "#" else 0 for c in row])
    row1, row2, col1, col2 = solve(grid)
    row1, row2, col1, col2 = row1 + 1, row2 + 1, col1 + 1, col2 + 1
    print(row1, row2)
    print(col1, col2)
