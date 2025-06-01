# 对于矩阵 grid 中的每个连续的 k x k 子矩阵，计算其中任意两个 不同值 之间的 最小绝对差 。

from itertools import pairwise
from typing import List


INF = int(1e18)


class Solution:
    def minAbsDiff(self, grid: List[List[int]], k: int) -> List[List[int]]:
        m, n = len(grid), len(grid[0])
        res = [[0] * (n - k + 1) for _ in range(m - k + 1)]
        for i in range(m - k + 1):
            subGrid = grid[i : i + k]
            for j in range(n - k + 1):
                rows = []
                for row in subGrid:
                    rows.extend(row[j : j + k])
                rows.sort()

                minDiff = INF
                for a, b in pairwise(rows):
                    if a < b:
                        minDiff = min(minDiff, b - a)
                if minDiff < INF:
                    res[i][j] = minDiff

        return res
