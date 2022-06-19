from heapq import nlargest
from typing import List


class XORMatrix:
    """二维前缀和模板(矩阵不可变)"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] ^ preSum[r][c + 1] ^ preSum[r + 1][c] ^ preSum[r][c]
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值::

        preSumMatrix.sumRegion(0, 0, 2, 2) # 左上角(0, 0)到右下角(2, 2)的值
        """
        return (
            self.preSum[r2 + 1][c2 + 1]
            ^ self.preSum[r2 + 1][c1]
            ^ self.preSum[r1][c2 + 1]
            ^ self.preSum[r1][c1]
        )


class Solution:
    def kthLargestValue(self, matrix: List[List[int]], k: int) -> int:
        M = XORMatrix(matrix)
        res = [M.sumRegion(0, 0, i, j) for i in range(len(matrix)) for j in range(len(matrix[0]))]
        return nlargest(k, res)[-1]
