from typing import List


# dp[i][j] 表示以 (i, j) 为右下角，且只包含 1 的正方形的边长最大值。
class Solution:
    def maximalSquare(self, matrix: List[List[str]]) -> int:
        ROW, COL = len(matrix), len(matrix[0])
        res = 0  # 最大正方形的边长
        dp = [[0] * COL for _ in range(ROW)]
        for row in range(ROW):
            for col in range(COL):
                if matrix[row][col] == "1":
                    if row == 0 or col == 0:
                        dp[row][col] = 1
                    else:
                        dp[row][col] = (
                            min(dp[row - 1][col - 1], dp[row - 1][col], dp[row][col - 1]) + 1
                        )
                    res = max(res, dp[row][col])
        return res * res


assert (
    Solution().maximalSquare(
        matrix=[
            ["1", "0", "1", "0", "0"],
            ["1", "0", "1", "1", "1"],
            ["1", "1", "1", "1", "1"],
            ["1", "0", "0", "1", "0"],
        ]
    )
    == 4
)
