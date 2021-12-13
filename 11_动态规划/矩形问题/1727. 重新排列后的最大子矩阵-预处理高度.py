from typing import List

# 你可以将 matrix 中的 列 按任意顺序重新排列。
# 请你返回最优方案下将 matrix 重新排列后，全是 1 的子矩阵面积。


# 存每行的高，然后可以参考leetcode 85的方法求每行最大矩形
# 因为列可以随意变换，这题直接对每行的高排序，然后求最大矩形就好
class Solution:
    def largestSubmatrix(self, matrix: List[List[int]]) -> int:
        row, col = len(matrix), len(matrix[0])

        def maxSquare(row: List[int]):
            res = 0
            for i, height in enumerate(row):
                res = max(res, height * (i + 1))
            return res

        #  最大矩形,预处理高度
        for r in range(row):
            for c in range(col):
                if matrix[r][c] == 1 and r > 0:
                    matrix[r][c] = matrix[r - 1][c] + 1

        matrix = [sorted(row, reverse=True) for row in matrix]
        return max(map(maxSquare, matrix))


print(Solution().largestSubmatrix(matrix=[[0, 0, 1], [1, 1, 1], [1, 0, 1]]))
# 输出：4
# 解释：你可以按照上图方式重新排列矩阵的每一列。
# 最大的全 1 子矩阵是上图中加粗的部分，面积为 4 。
