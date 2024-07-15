# 1444. 切披萨的方案数
# https://leetcode.cn/problems/number-of-ways-of-cutting-a-pizza/description/
# 切披萨的每一刀，先要选择是向垂直还是水平方向切
# 你需要切披萨 k-1 次，得到 k 块披萨并送给别人。
# !如果垂直地切披萨，那么需要把左边的部分送给一个人，如果水平地切，那么需要把上面的部分送给一个人.
# 请你返回确保每一块披萨包含 至少 一个苹果的切披萨方案数。由于答案可能是个很大的数字，请你返回它对 10^9 + 7 取余的结果。
#
# 1 <= rows, cols <= 50
# 1 <= k <= 10

from typing import List
from functools import lru_cache


MOD = int(1e9 + 7)


class Solution:
    def ways(self, pizza: List[str], k: int) -> int:
        """
        你需要切披萨 k-1 次，得到 k 块披萨并送给别人。
        请你返回确保每一块披萨包含 至少 一个苹果的切披萨方案数
        """

        @lru_cache(None)
        def dfs(row: int, col: int, remain: int) -> int:
            """左上角(row,col)到右下角(ROW-1,COL-1)的披萨,还剩remain刀的切法"""
            if remain <= 0:
                return 1 if P.sumRegion(row, ROW - 1, col, COL - 1) > 0 else 0

            res = 0

            # !横着切
            for r in range(row, ROW):
                if P.sumRegion(row, r, col, COL - 1) > 0:
                    res += dfs(r + 1, col, remain - 1)
            # !竖着切
            for c in range(col, COL):
                if P.sumRegion(row, ROW - 1, col, c) > 0:
                    res += dfs(row, c + 1, remain - 1)

            return res % MOD

        P = PreSum2DDense([[1 if v == "A" else 0 for v in list(row)] for row in pizza])
        ROW, COL = len(pizza), len(pizza[0])
        return dfs(0, 0, k - 1)


class PreSum2DDense:
    """二维前缀和模板(矩阵不可变)"""

    __slots__ = "_preSum"

    def __init__(self, mat: List[List[int]]):
        ROW, COL = len(mat), len(mat[0])
        preSum = [[0] * (COL + 1) for _ in range(ROW + 1)]
        for r in range(ROW):
            tmpSum0, tmpSum1 = preSum[r], preSum[r + 1]
            tmpM = mat[r]
            for c in range(COL):
                tmpSum1[c + 1] = tmpM[c] + tmpSum0[c + 1] + tmpSum1[c] - tmpSum0[c]
        self._preSum = preSum

    def sumRegion(self, x1: int, x2: int, y1: int, y2: int) -> int:
        """查询sum(A[x1:x2+1, y1:y2+1])的值(包含边界)."""
        if x1 > x2 or y1 > y2:
            return 0
        return (
            self._preSum[x2 + 1][y2 + 1]
            - self._preSum[x2 + 1][y1]
            - self._preSum[x1][y2 + 1]
            + self._preSum[x1][y1]
        )


if __name__ == "__main__":
    # assert Solution().ways(pizza=["A..", "AAA", "..."], k=3) == 3
    print(Solution().ways(pizza=[".A..A", "A.A..", "A.AA.", "AAAA.", "A.AA."], k=5))
    assert Solution().ways(pizza=[".A..A", "A.A..", "A.AA.", "AAAA.", "A.AA."], k=5) == 153
