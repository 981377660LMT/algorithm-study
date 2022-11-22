# https://atcoder.jp/contests/abc175/tasks/abc175_e

# !输入 n m (1≤n,m≤3000) k(≤min(2e5,r*c))，表示一个 n*m 的网格，和网格中的 k 个物品。
# 接下来 k 行，每行三个数 x y v(≤1e9) 表示物品的行号、列号和价值（行列号从 1 开始）。
# 每个网格至多有一个物品。

# 你从 (1,1) 出发走到 (n,m)，每步只能向下或向右。
# 经过物品时，你可以选或不选，且每行至多可以选三个物品。
# 输出你选到的物品的价值和的最大值。

from functools import lru_cache
from typing import List
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


def solve(grid: List[List[int]]) -> int:
    @lru_cache(None)
    def dfs(row: int, col: int, count: int) -> int:
        if row == ROW - 1 and col == COL - 1:
            return 0

        res = 0
        if row + 1 < ROW:
            res = max(res, dfs(row + 1, col, 0))
            if grid[row + 1][col] > 0:
                res = max(res, grid[row + 1][col] + dfs(row + 1, col, 1))
        if col + 1 < COL:
            res = max(res, dfs(row, col + 1, count))
            if grid[row][col + 1] > 0 and count <= 2:
                res = max(res, grid[row][col + 1] + dfs(row, col + 1, count + 1))

        return res

    res = dfs(0, 0, 0)
    if grid[0][0] > 0:  # 有物品
        res = max(res, grid[0][0] + dfs(0, 0, 1))
    return res


if __name__ == "__main__":
    ROW, COL, K = map(int, input().split())
    grid = [[0] * COL for _ in range(ROW)]
    for _ in range(K):
        row, col, score = map(int, input().split())
        row, col = row - 1, col - 1
        grid[row][col] = score
    print(solve(grid))
