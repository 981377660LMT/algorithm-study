from typing import List

# 一个「角矩形」是由`四个`不同的在网格上的 1 形成的轴对称的矩形。注意只有角的位置才需要为 1。并且，4 个 1 需要是不同的
# 网格 grid 中行和列的数目范围为 [1, 200]。


class Solution:
    def countCornerRectangles(self, grid: List[List[int]]) -> int:
        """相邻行间进行dp"""
        row, col = len(grid), len(grid[0])
        res = 0

        dp = [[0] * col for _ in range(col)]  # dp[i][j] 表示以 i, j 为左、右下角的角矩形的数量
        for r in range(row):
            for leftC in range(col):
                if grid[r][leftC] == 1:
                    for rightC in range(leftC + 1, col):
                        if grid[r][rightC] == 1:
                            res += dp[leftC][rightC]
                            dp[leftC][rightC] += 1

        return res


print(
    Solution().countCornerRectangles(
        grid=[[1, 0, 0, 1, 0], [0, 0, 1, 0, 1], [0, 0, 0, 1, 0], [1, 0, 1, 0, 1]]
    )
)
# 输出：1
# 解释：只有一个角矩形，角的位置为 grid[1][2], grid[1][4], grid[3][2], grid[3][4]。
