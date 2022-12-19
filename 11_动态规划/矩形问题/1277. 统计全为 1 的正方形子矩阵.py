from typing import List

# 给你一个 m * n 的矩阵，矩阵中的元素不是 0 就是 1，
# 请你统计并返回其中完全由 1 组成的 正方形 子矩阵的个数。

# dp[i][j] 表示以 (i, j) 为右下角，且只包含 1 的正方形的边长最大值。
# 统计时，计算以 (i, j) 为右下角的正方形的个数，即 dp[i][j]。


class Solution:
    def countSquares(self, A: List[List[int]]) -> int:
        ROW, COL = len(A), len(A[0])
        dp = [[0] * COL for _ in range(ROW)]
        res = 0
        for row in range(ROW):
            for col in range(COL):
                if A[row][col] == 1:
                    if row == 0 or col == 0:
                        dp[row][col] = 1
                    else:
                        dp[row][col] = (
                            min(dp[row - 1][col - 1], dp[row - 1][col], dp[row][col - 1]) + 1
                        )
                    res += dp[row][col]
        return res


assert Solution().countSquares(A=[[0, 1, 1, 1], [1, 1, 1, 1], [0, 1, 1, 1]]) == 15
# 解释：
# 边长为 1 的正方形有 10 个。
# 边长为 2 的正方形有 4 个。
# 边长为 3 的正方形有 1 个。
# 正方形的总数 = 10 + 4 + 1 = 15.
