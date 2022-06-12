from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minPathCost(self, grid: List[List[int]], moveCost: List[List[int]]) -> int:
        # 有路径经过的单元格的 值之和 加上 所有移动的 代价之和
        ROW, COL = len(grid), len(grid[0])
        dp = [(num, num) for num in grid[0]]  # cost value
        for r in range(1, ROW):
            ndp = [(int(1e20), -1)] * COL
            for c, num in enumerate(grid[r]):
                for pre in range(COL):
                    preCost, preNum = dp[pre]
                    cand = preCost + moveCost[preNum][c] + num
                    if cand < ndp[c][0]:
                        ndp[c] = (cand, num)
            dp = ndp
        return min([num for num, _ in dp])

