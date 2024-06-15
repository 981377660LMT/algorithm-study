# G - AtCoder Tour
# https://atcoder.jp/contests/abc358/tasks/abc358_g
# 给定一个n*m的网格，每个格子有一个数字.
# 从(sx,sy)出发，每次可以不动，或者向上下左右移动一格.
# 求移动k次后的得分最大值.
# n,m<=50,k<=1e9.
#
# 非常像矩阵快速幂，但是转移矩阵太大了，不可行.
# !考虑dp，最后一定是移动到某个格子后一直不动.
# dp[i][j][k]表示从(sx,sy)出发，移动k次到达(i,j)的最大得分.


from typing import List


INF = int(4e18)
DIR4 = [(0, 1), (0, -1), (1, 0), (-1, 0)]


def max2(a: int, b: int) -> int:
    return a if a > b else b


def longestPathInGrid(grid: List[List[int]], sx: int, sy: int, k: int) -> int:
    ROW, COL = len(grid), len(grid[0])
    dp = [[-INF] * COL for _ in range(ROW)]
    dp[sx][sy] = 0

    res = -INF
    for round in range(ROW * COL):
        if round > k:
            break
        for i in range(ROW):
            for j in range(COL):
                res = max2(res, dp[i][j] + grid[i][j] * (k - round))

        ndp = [list(row) for row in dp]
        for i in range(ROW):
            for j in range(COL):
                for di, dj in DIR4:
                    ni, nj = i + di, j + dj
                    if ni < 0 or ni >= ROW or nj < 0 or nj >= COL:
                        continue
                    ndp[ni][nj] = max2(ndp[ni][nj], dp[i][j] + grid[ni][nj])

        dp = ndp

    return res


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    sx, sy = map(int, input().split())
    sx, sy = sx - 1, sy - 1
    grid = [list(map(int, input().split())) for _ in range(n)]
    print(longestPathInGrid(grid, sx, sy, k))
