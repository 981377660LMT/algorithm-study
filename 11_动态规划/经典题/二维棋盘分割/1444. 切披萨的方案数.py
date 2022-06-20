from typing import List
from functools import lru_cache

# 1 <= rows, cols <= 50
# 1 <= k <= 10

# 切披萨的每一刀，先要选择是向垂直还是水平方向切
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
        if r1 > r2 or c1 > c2:
            return 0
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
        def dfs(sr: int, sc: int, k: int) -> int:
            """左上角(x, y)到右下角(row-1, col-1)能否再切k次"""
            if k == 0:
                return int(M.sumRegion(sr, sc, ROW - 1, COL - 1) > 0)

            res = 0

            # !下一块左上角的位置
            for r in range(sr + 1, ROW):
                if M.sumRegion(sr, sc, r - 1, COL - 1) > 0:
                    res += dfs(r, sc, k - 1)
                    res %= MOD

            for c in range(sc + 1, COL):
                if M.sumRegion(sr, sc, ROW - 1, c - 1) > 0:
                    res += dfs(sr, c, k - 1)
                    res %= MOD

            return res

        M = PreSumMatrix([list(row) for row in pizza])
        ROW, COL = len(pizza), len(pizza[0])
        return dfs(0, 0, k - 1)


print(Solution().ways(pizza=["A..", "AAA", "..."], k=3))
# 输出：3
# 解释：上图展示了三种切披萨的方案。注意每一块披萨都至少包含一个苹果。
