from typing import List
from scipy.optimize import linear_sum_assignment


# scipy调库解法
class Solution:
    def maxCompatibilitySum(self, students: List[List[int]], mentors: List[List[int]]) -> int:
        n, m = (len(students), len(students[0]))
        costMatrix = [[0] * n for _ in range(n)]
        for i in range(n):
            for j in range(n):
                for k in range(m):
                    costMatrix[i][j] += int(students[i][k] == mentors[j][k])

        rowIndex, colIndex = linear_sum_assignment(costMatrix, maximize=True)
        return sum(costMatrix[row][col] for row, col in zip(rowIndex, colIndex))


print(
    Solution().maxCompatibilitySum(
        students=[[1, 1, 0], [1, 0, 1], [0, 0, 1]], mentors=[[1, 0, 0], [0, 0, 1], [1, 1, 0]]
    )
)
# 解释：按下述方式分配学生和导师：
# - 学生 0 分配给导师 2 ，兼容性评分为 3 。
# - 学生 1 分配给导师 0 ，兼容性评分为 2 。
# - 学生 2 分配给导师 1 ，兼容性评分为 3 。
# 最大兼容性评分和为 3 + 2 + 3 = 8 。
