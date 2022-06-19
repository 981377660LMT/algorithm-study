from typing import List
from functools import lru_cache

# 1 <= rows, cols <= 50
# 1 <= k <= 10

# 你需要切披萨 k-1 次，得到 k 块披萨并送给别人。
# 请你返回确保每一块披萨包含 至少 一个苹果的切披萨方案数。由于答案可能是个很大的数字，请你返回它对 10^9 + 7 取余的结果。


MOD = int(1e9 + 7)


class PreSumMatrix:
    def __init__(self, A: List[List[str]]):
        m, n = len(A), len(A[0])

        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = (
                    int(A[r][c] == 'A') + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
                )

        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


class Solution:
    def ways(self, pizza: List[str], k: int) -> int:
        """你需要切披萨 k-1 次，得到 k 块披萨并送给别人。
        
        请你返回确保每一块披萨包含 至少 一个苹果的切披萨方案数
        """

        @lru_cache(None)
        def dfs(x, y, k):
            """左上角(x, y)到右下角(row-1, col-1)能否再切k次"""
            if not k:
                return prefix[r][c] - prefix[x][c] - prefix[r][y] + prefix[x][y] > 0
            res = 0
            for i in range(x + 1, r):
                if prefix[i][c] - prefix[x][c] - prefix[i][y] + prefix[x][y] > 0:
                    res += dfs(i, y, k - 1)
            for j in range(y + 1, c):
                if prefix[r][j] - prefix[x][j] - prefix[r][y] + prefix[x][y] > 0:
                    res += dfs(x, j, k - 1)
            return res

        M = PreSumMatrix([list(row) for row in pizza])
        return dfs(0, 0, k - 1)


print(Solution().ways(pizza=["A..", "AAA", "..."], k=3))
# 输出：3
# 解释：上图展示了三种切披萨的方案。注意每一块披萨都至少包含一个苹果。
