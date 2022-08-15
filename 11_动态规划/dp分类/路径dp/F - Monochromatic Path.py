"""
可以用某个花费任意翻转行列
求左上角到右下角以相同颜色连通的最小花费


dp[row][col][rowFlip][colFlip]
"""

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

ROW, COL = map(int, input().split())
rowCost = list(map(int, input().split()))  # 翻转行的代价
colCost = list(map(int, input().split()))  # 翻转列的代价
grid = []
for _ in range(ROW):
    grid.append([int(char) for char in input()])

# !同色连通左上角到右下角的最小花费

# !TLE (atc数据大的题不要用dfs写dp)
# def dfs(row: int, col: int, rowFlip: bool, colFlip: bool) -> int:
#     """行列、行是否翻转、列是否翻转"""
#     if row == ROW - 1 and col == COL - 1:
#         return 0

#     hash_ = row * COL * 4 + col * 4 + rowFlip * 2 + colFlip
#     if ~memo[hash_]:
#         return memo[hash_]

#     res = INF
#     # 向下走 可能需要翻转行
#     if row + 1 < ROW:
#         if grid[row + 1][col] == (grid[row][col] ^ rowFlip):
#             res = min(res, dfs(row + 1, col, False, colFlip))
#         else:
#             res = min(res, dfs(row + 1, col, True, colFlip) + rowCost[row + 1])

#     # 向右走 可能需要翻转列
#     if col + 1 < COL:
#         if grid[row][col + 1] == (grid[row][col] ^ colFlip):
#             res = min(res, dfs(row, col + 1, rowFlip, False))
#         else:
#             res = min(res, dfs(row, col + 1, rowFlip, True) + colCost[col + 1])

#     memo[hash_] = res
#     return res


# memo = [-1] * ((ROW + 1) * (COL + 1) * 4)
# # 注意起始位置的翻转
# print(
#     min(
#         dfs(0, 0, False, False),
#         dfs(0, 0, True, False) + rowCost[0],
#         dfs(0, 0, False, True) + colCost[0],
#         dfs(0, 0, True, True) + rowCost[0] + colCost[0],
#     )
# )


# !dp[行翻转][列翻转][坐标]
dp = [[[INF] * ROW * COL for _ in range(2)] for _ in range(2)]
dp[0][0][0] = 0
dp[1][0][0] = rowCost[0]
dp[0][1][0] = colCost[0]
dp[1][1][0] = rowCost[0] + colCost[0]

for r in range(ROW):
    for c in range(COL):
        pos = r * COL + c
        for rowFlip in range(2):
            for colFlip in range(2):
                if r + 1 < ROW:
                    if grid[r + 1][c] == (grid[r][c] ^ rowFlip):
                        dp[0][colFlip][pos + COL] = min(
                            dp[0][colFlip][pos + COL], dp[rowFlip][colFlip][pos]
                        )
                    else:
                        dp[1][colFlip][pos + COL] = min(
                            dp[1][colFlip][pos + COL], dp[rowFlip][colFlip][pos] + rowCost[r + 1]
                        )
                if c + 1 < COL:
                    if grid[r][c + 1] == (grid[r][c] ^ colFlip):
                        dp[rowFlip][0][pos + 1] = min(
                            dp[rowFlip][0][pos + 1], dp[rowFlip][colFlip][pos]
                        )
                    else:
                        dp[rowFlip][1][pos + 1] = min(
                            dp[rowFlip][1][pos + 1], dp[rowFlip][colFlip][pos] + colCost[c + 1]
                        )

print(min(dp[1][1][-1], dp[1][0][-1], dp[0][1][-1], dp[0][0][-1]))
