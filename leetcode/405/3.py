from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二维字符矩阵 grid，其中 grid[i][j] 可能是 'X'、'Y' 或 '.'，返回满足以下条件的子矩阵数量：


# 包含 grid[0][0]
# 'X' 和 'Y' 的频数相等。
# 至少包含一个 'X'。


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


class Solution:
    def numberOfSubmatrices(self, grid: List[List[str]]) -> int:
        xGrid = [[1 if c == "X" else 0 for c in row] for row in grid]
        yGrid = [[1 if c == "Y" else 0 for c in row] for row in grid]
        xSum = PreSum2DDense(xGrid)
        ySum = PreSum2DDense(yGrid)
        res = 0
        ROW, COL = len(grid), len(grid[0])
        for x2 in range(ROW):
            for y2 in range(COL):
                sum1 = xSum.sumRegion(0, x2, 0, y2)
                sum2 = ySum.sumRegion(0, x2, 0, y2)
                if sum1 == sum2 and sum1 > 0:
                    res += 1
        return res
