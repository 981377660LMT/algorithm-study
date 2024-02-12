from typing import List


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

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self._preSum[r2 + 1][c2 + 1]
            - self._preSum[r2 + 1][c1]
            - self._preSum[r1][c2 + 1]
            + self._preSum[r1][c1]
        )
