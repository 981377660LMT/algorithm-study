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
