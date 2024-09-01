# https://leetcode.cn/problems/select-cells-in-grid-with-maximum-score/description/
# 3276. 选择矩阵中单元格的最大得分
# 你需要从矩阵中选择 一个或多个 单元格，选中的单元格应满足以下条件：
# !所选单元格中的任意两个单元格都不会处于矩阵的 同一行。
# !所选单元格的值 互不相同。
# !你的得分为所选单元格值的总和。
# 返回你能获得的 最大 得分。
#
# 1 <= grid.length, grid[i].length <= 10
# 1 <= grid[i][j] <= 100
#
# 我们建立一个二分图，左部点为 100 个数字，右部点为 n 行，对于数字 i，如果其在第 j 行出现，则向第 j 个右部点连流量为 1 的边。
# 随后，源点向每个数字连流量为 1 的边，代表每个数字只能选一次，同时计入等于数字的费用；
# 每行向汇点连流量为 1 的边，表示每行只能选一个。
# 根据这个模型直接跑最大费用最大流，或者写二分图最大权匹配，
# 又或者转换成模拟费用流模型，变成经典的增广路求二分图最大匹配都是可以过的

from typing import List
from scipy.optimize import linear_sum_assignment


class Solution:
    def maxScore(self, grid: List[List[int]]) -> int:
        max_ = max(max(row) for row in grid)
        costMatrix = [[0] * len(grid) for _ in range(max_ + 1)]
        for r, row in enumerate(grid):
            for _, val in enumerate(row):
                costMatrix[val][r] = val
        rowIndex, colIndex = linear_sum_assignment(costMatrix, maximize=True)
        return sum(costMatrix[row][col] for row, col in zip(rowIndex, colIndex))
