from functools import lru_cache
from typing import List


class PreSumMatrix:
    """二维前缀和矩阵"""

    def __init__(self, A: List[List[int]]):
        m, n = len(A), len(A[0])
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(m):
            for c in range(n):
                preSum[r + 1][c + 1] = A[r][c] + preSum[r][c + 1] + preSum[r + 1][c] - preSum[r][c]
        self.preSum = preSum

    def sumRegion(self, r1: int, c1: int, r2: int, c2: int) -> int:
        """查询sum(A[r1:r2+1, c1:c2+1])的值"""
        return (
            self.preSum[r2 + 1][c2 + 1]
            - self.preSum[r2 + 1][c1]
            - self.preSum[r1][c2 + 1]
            + self.preSum[r1][c1]
        )


# 启示：与类无关的函数可以丢到类外面预处理，加一层记忆化
@lru_cache(None)
def getFactorCount(num: int, factor: int) -> int:
    """预处理因子个数"""
    res = 0
    while num % factor == 0:
        num = num // factor
        res += 1
    return res


# 二维前缀和，把计算前缀和单独提取出一个类，代码还是很清晰的
class Solution:
    def maxTrailingZeros(self, grid: List[List[int]]) -> int:
        """请你从 grid 中找出一条乘积中尾随零数目最多的转角路径，并返回该路径中尾随零的数目。"""
        row, col = len(grid), len(grid[0])
        grid2 = [[getFactorCount(num, 2) for num in row] for row in grid]
        grid5 = [[getFactorCount(num, 5) for num in row] for row in grid]
        preSumMatrix2 = PreSumMatrix(grid2)
        preSumMatrix5 = PreSumMatrix(grid5)

        res = 0
        for r in range(row):
            for c in range(col):
                up2 = preSumMatrix2.sumRegion(0, c, r, c)
                right2 = preSumMatrix2.sumRegion(r, c, r, col - 1)
                left2 = preSumMatrix2.sumRegion(r, 0, r, c)
                down2 = preSumMatrix2.sumRegion(r, c, row - 1, c)
                up5 = preSumMatrix5.sumRegion(0, c, r, c)
                right5 = preSumMatrix5.sumRegion(r, c, r, col - 1)
                left5 = preSumMatrix5.sumRegion(r, 0, r, c)
                down5 = preSumMatrix5.sumRegion(r, c, row - 1, c)

                cand1 = min(up5 + left5 - grid5[r][c], up2 + left2 - grid2[r][c])
                cand2 = min(right5 + down5 - grid5[r][c], right2 + down2 - grid2[r][c])
                cand3 = min(up5 + right5 - grid5[r][c], up2 + right2 - grid2[r][c])
                cand4 = min(left5 + down5 - grid5[r][c], left2 + down2 - grid2[r][c])

                res = max(res, cand1, cand2, cand3, cand4)

        return res


print(
    Solution().maxTrailingZeros(
        grid=[
            [23, 17, 15, 3, 20],
            [8, 1, 20, 27, 11],
            [9, 4, 6, 2, 21],
            [40, 9, 1, 10, 6],
            [22, 7, 4, 5, 3],
        ]
    )
)

print(
    Solution().maxTrailingZeros(
        [
            [899, 727, 165, 249, 531, 300, 542, 890],
            [981, 587, 565, 943, 875, 498, 582, 672],
            [106, 902, 524, 725, 699, 778, 365, 220],
        ]
    )
)

# 输出：
# 6
# 预期：
# 5
