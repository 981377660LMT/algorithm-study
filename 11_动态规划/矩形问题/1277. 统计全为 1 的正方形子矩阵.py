from typing import List

# 221. 最大正方形
# Note this is almost identical to #221 just change the max there to + here
# dp[i][j] means the size of biggest square with A[i][j] as bottom-right corner.
# dp[i][j] also means the number of squares with A[i][j] as bottom-right corner.


class Solution:
    def countSquares(self, A: List[List[int]]) -> int:
        for i in range(1, len(A)):
            for j in range(1, len(A[0])):
                A[i][j] *= min(A[i - 1][j], A[i][j - 1], A[i - 1][j - 1]) + 1
        return sum(sum(row) for row in A)


print(Solution().countSquares(A=[[0, 1, 1, 1], [1, 1, 1, 1], [0, 1, 1, 1]]))
# 解释：
# 边长为 1 的正方形有 10 个。
# 边长为 2 的正方形有 4 个。
# 边长为 3 的正方形有 1 个。
# 正方形的总数 = 10 + 4 + 1 = 15.

