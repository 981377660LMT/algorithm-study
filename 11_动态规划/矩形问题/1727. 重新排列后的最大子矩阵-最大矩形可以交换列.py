from typing import List

# 你可以将 matrix 中的 列 按任意顺序重新排列。
# 请你返回最优方案下将 matrix 重新排列后，全是 1 的子矩阵面积。
# ROW*COL<=1e5

# !`最大矩形可以交换列的版本=>预处理+排序`
# 存每行的高，然后可以参考leetcode 85的方法求每行最大矩形
# 因为列可以随意变换，这题直接对每行的高排序，然后求最大矩形就好


class Solution:
    def largestSubmatrix(self, matrix: List[List[int]]) -> int:
        ROW, COL = len(matrix), len(matrix[0])
        up = [[0] * COL for _ in range(ROW)]  # 预处理高度

        for r in range(ROW):
            for c in range(COL):
                if matrix[r][c] == 1:
                    up[r][c] = up[r - 1][c] + 1 if r else 1

        res = 0
        for row in up:
            row.sort(reverse=True)
            cand = max(h * (i + 1) for i, h in enumerate(row))
            res = max(res, cand)
        return res


print(Solution().largestSubmatrix(matrix=[[0, 0, 1], [1, 1, 1], [1, 0, 1]]))
# 输出：4
# 解释：你可以按照上图方式重新排列矩阵的每一列。
# 最大的全 1 子矩阵是上图中加粗的部分，面积为 4 。
