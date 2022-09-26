"""最大子矩阵"""

from typing import List, Tuple


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


def maxSubArray(nums: List[int]) -> Tuple[int, Tuple[int, int]]:
    """最大子数组和,返回数组

    dp,需要记录左端点:只取当前还是取前面
    如果前面的和小于0,那么就舍弃前面的一截,并将左端点移到当前位置
    """

    if len(nums) == 1:
        return nums[0], (0, 0)

    maxSum, curSum = -int(1e20), 0
    preLeft = 0
    resLeft, resRight = 0, 0
    for i, num in enumerate(nums):
        if curSum < 0:
            curSum = num
            preLeft = i
        else:
            curSum += num

        if curSum > maxSum:
            maxSum = curSum
            resLeft = preLeft
            resRight = i

    return maxSum, (resLeft, resRight)


# -100 <= matrix[i][j] <= 100
# 1 <= m, n <= 100
class Solution:
    def getMaxMatrix(self, matrix: List[List[int]]) -> List[int]:
        """给你一个 m x n 的矩阵 matrix 和一个整数 k ，找出元素总和最大的子矩阵。"""
        ROW, COL = len(matrix), len(matrix[0])
        P = PreSumMatrix(matrix)
        res = [0, 0, 0, 0]
        maxSum = -int(1e20)

        for r1 in range(ROW):
            for r2 in range(r1, ROW):
                nums = [P.sumRegion(r1, c, r2, c) for c in range(COL)]
                curMax, (left, right) = maxSubArray(nums)
                if curMax >= maxSum:
                    maxSum = curMax
                    res = [r1, left, r2, right]
        return res


if __name__ == "__main__":
    matrix = [[9, -8, 1, 3, -2], [-3, 7, 6, -2, 4], [6, -4, -4, 8, -7]]
    assert Solution().getMaxMatrix(matrix) == [0, 0, 2, 3]
