from typing import List


class Solution:
    def minPathCost(self, grid: List[List[int]], moveCost: List[List[int]]) -> int:
        # 有路径经过的单元格的 值之和 加上 所有移动的 代价之和
        ROW, COL = len(grid), len(grid[0])
        dp = [num for num in grid[0]]
        for r in range(1, ROW):
            ndp = [int(1e20)] * COL
            for curC, num in enumerate(grid[r]):
                for preC in range(COL):
                    preCost = dp[preC]
                    ndp[curC] = min(ndp[curC], preCost + moveCost[grid[r - 1][preC]][curC] + num)
            dp = ndp
        return min(dp)


print(
    Solution().minPathCost(
        grid=[[5, 3], [4, 0], [2, 1]], moveCost=[[9, 8], [1, 5], [10, 12], [18, 6], [2, 4], [14, 3]]
    )
)

