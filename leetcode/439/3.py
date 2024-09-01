from heapq import heapify
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个由正整数构成的二维矩阵 grid。

# 你需要从矩阵中选择 一个或多个 单元格，选中的单元格应满足以下条件：

# 所选单元格中的任意两个单元格都不会处于矩阵的 同一行。
# 所选单元格的值 互不相同。
# 你的得分为所选单元格值的总和。


# 返回你能获得的 最大 得分。


class Solution:
    def maxScore(self, grid: List[List[int]]) -> int:
        grid = [sorted(set(row), reverse=True) for row in grid]
        ROW = len(grid)
        pq = [(grid[r][0], r, 0) for r in range(ROW)]
        heapify(pq)
        max_ = max(item[0] for item in pq)
        while True:
            min_, row, col = heappop(pq)
            if max_ - min_ < rightRes - leftRes:
                leftRes, rightRes = min_, max_
            if col == len(nums[row]) - 1:
                break
            max_ = max(max_, nums[row][col + 1])
            heappush(pq, (nums[row][col + 1], row, col + 1))
        return [leftRes, rightRes]
