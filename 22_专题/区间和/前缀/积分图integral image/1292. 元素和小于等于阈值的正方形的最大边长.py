from typing import List

# 请你返回元素总和小于或等于阈值的正方形区域的最大边长
class Solution:
    def maxSideLength(self, mat: List[List[int]], threshold: int) -> int:
        m, n = len(mat), len(mat[0])
        preSum = [[0] * (n + 1) for _ in range(m + 1)]
        for r in range(1, m + 1):
            for c in range(1, n + 1):
                preSum[r][c] = (
                    preSum[r - 1][c] + preSum[r][c - 1] - preSum[r - 1][c - 1] + mat[r - 1][c - 1]
                )

        def getSum(x1, y1, x2, y2):
            return preSum[x2][y2] - preSum[x1 - 1][y2] - preSum[x2][y1 - 1] + preSum[x1 - 1][y1 - 1]

        def check(mid) -> bool:
            for i in range(1, m - mid + 2):
                for j in range(1, n - mid + 2):
                    if getSum(i, j, i + mid - 1, j + mid - 1) <= threshold:
                        print(i, j, mid)
                        return True
            return False

        left, right = 1, min(m, n)
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
