from collections import Counter
from typing import List
from sortedcontainers import SortedList


class PreSumMatrix:
    """二维前缀和模板(矩阵不可变)"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])

        # 前缀和数组
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
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


# -100 <= matrix[i][j] <= 100
# 1 <= m, n <= 100


class Solution:
    def numSubmatrixSumTarget(self, matrix: List[List[int]], k: int) -> int:
        """给你一个 m x n 的矩阵 matrix 和一个整数 k ，返回元素总和等于目标值的非空子矩阵的数量。。
        
        暴力枚举需要O(m^2*n^2)
        优化:枚举上边界和下边界O(n^2),就变成了一维的找一个子数组使得和最接近k
        把前缀和记录到有序集合里，然后二分寻找 O(m*n*min(m,n)log(min(m,n)))
        如果行数远大于列数，可以先转置矩阵
        """
        ROW, COL = len(matrix), len(matrix[0])
        P = PreSumMatrix(matrix)
        res = 0

        for r1 in range(ROW):
            for r2 in range(r1, ROW):
                counter, curSum = Counter({0: 1}), 0
                for c in range(COL):
                    curSum += P.sumRegion(r1, c, r2, c)
                    res += counter[curSum - k]
                    counter[curSum] += 1
        return res
