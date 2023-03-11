from typing import List


class M:
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


# 请你返回元素总和小于或等于阈值的正方形区域的最大边长，元素都是非负数
class Solution:
    def maxSideLength(self, mat: List[List[int]], threshold: int) -> int:
        def check(mid: int) -> bool:
            for r in range(row):
                for c in range(col):
                    if r + mid - 1 < row and c + mid - 1 < col:
                        if p.sumRegion(r, c, r + mid - 1, c + mid - 1) <= threshold:
                            return True
            return False

        row, col = len(mat), len(mat[0])
        p = M(mat)
        left, right = 0, min(row, col)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1

        return right


print(
    Solution().maxSideLength(
        mat=[[1, 1, 3, 2, 4, 3, 2], [1, 1, 3, 2, 4, 3, 2], [1, 1, 3, 2, 4, 3, 2]], threshold=4
    )
)

# 前缀和：p[i-1][j]+p[i][j-1]-p[i-1][j-1]+mat[i-1][j-1]
# 矩形和：p[x2][y2] - p[x1-1][y2] - p[x2][y1-1] + p[x1-1][y1-1]
