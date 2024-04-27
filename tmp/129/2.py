from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二维 boolean 矩阵 grid 。

# 请你返回使用 grid 中的 3 个元素可以构建的 直角三角形 数目，且满足 3 个元素值 都 为 1 。

# 注意：


# 如果 grid 中 3 个元素满足：一个元素与另一个元素在 同一行，同时与第三个元素在 同一列 ，那么这 3 个元素称为一个 直角三角形 。这 3 个元素互相之间不需要相邻。


class PreSum2D:
    """二维前缀和模板."""

    __slots__ = "preSum"

    def __init__(self, mat: List[List[int]]):
        ROW, COL = len(mat), len(mat[0])
        preSum = [[0] * (COL + 1) for _ in range(ROW + 1)]
        for r in range(ROW):
            tmpSum0, tmpSum1 = preSum[r], preSum[r + 1]
            tmpM = mat[r]
            for c in range(COL):
                tmpSum1[c + 1] = tmpM[c] + tmpSum0[c + 1] + tmpSum1[c] - tmpSum0[c]
        self.preSum = preSum

    def queryRange(self, r1: int, c1: int, r2: int, c2: int) -> int:
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
    def numberOfRightTriangles(self, grid: List[List[int]]) -> int:
        S = PreSum2D(grid)
        ROW, COL = len(grid), len(grid[0])

        # 枚举直角顶点
        res = 0
        for i in range(len(grid)):
            for j in range(len(grid[0])):
                if grid[i][j] == 1:
                    upSum = S.queryRange(0, j, i, j) - 1
                    downSum = S.queryRange(i, j, ROW - 1, j) - 1
                    leftSum = S.queryRange(i, 0, i, j) - 1
                    rightSum = S.queryRange(i, j, i, COL - 1) - 1
                    res += (
                        upSum * leftSum + upSum * rightSum + downSum * leftSum + downSum * rightSum
                    )

        return res
