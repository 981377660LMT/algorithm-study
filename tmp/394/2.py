from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个大小为 m x n 的二维矩形 grid 。每次 操作 中，你可以将 任一 格子的值修改为 任意 非负整数。完成所有操作后，你需要确保每个格子 grid[i][j] 的值满足：


# 如果下面相邻格子存在的话，它们的值相等，也就是 grid[i][j] == grid[i + 1][j]（如果存在）。
# 如果右边相邻格子存在的话，它们的值不相等，也就是 grid[i][j] != grid[i][j + 1]（如果存在）。
# 请你返回需要的 最少 操作数目。


# !列相等，相邻列不等


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumOperations(self, grid: List[List[int]]) -> int:
        ROW, COL = len(grid), len(grid[0])
        colCounter = [defaultdict(int) for _ in range(COL)]
        for i, row in enumerate(grid):
            for j, num in enumerate(row):
                colCounter[j][num] += 1

        @lru_cache(None)
        def dfs(index: int, pre: int) -> int:
            if index == COL:
                return 0
            res = INF
            counter = colCounter[index]
            for cur in range(12):
                if cur == pre:
                    continue
                diff = ROW - counter[cur]
                res = min2(res, diff + dfs(index + 1, cur))
            return res

        res = dfs(0, -1)
        dfs.cache_clear()
        return res
