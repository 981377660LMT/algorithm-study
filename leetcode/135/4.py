from functools import lru_cache
from itertools import accumulate
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个大小为 n x n 的二维矩阵 grid ，一开始所有格子都是白色的。一次操作中，你可以选择任意下标为 (i, j) 的格子，并将第 j 列中从最上面到第 i 行所有格子改成黑色。

# 如果格子 (i, j) 为白色，且左边或者右边的格子至少一个格子为黑色，那么我们将 grid[i][j] 加到最后网格图的总分中去。


# 请你返回执行任意次操作以后，最终网格图的 最大 总分数。


# !黑色区域左右相邻的白色区域是可以加分的，求最大分数


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def maximumScore(self, grid: List[List[int]]) -> int:
        n = len(grid)
        ROW, COL = n, n
        colPreSum = [list(accumulate(col, initial=0)) for col in zip(*grid)]

        def cal(col: int, rowCount1: int, rowCount2: int) -> int:
            if col < 0 or col >= COL:
                return 0
            if rowCount1 >= rowCount2:
                return 0
            arr = colPreSum[col]
            return arr[rowCount2] - arr[rowCount1]

        @lru_cache(None)
        def dfs(index: int, preBlackRowCount1: int, preBlackRowCount2: int) -> int:
            if index == COL:
                return cal(index - 1, preBlackRowCount1, preBlackRowCount2)
            res = 0
            for i in range(ROW + 1):
                pre = cal(index - 1, preBlackRowCount1, max2(i, preBlackRowCount2))  # 前一列白色得分
                res = max2(res, pre + dfs(index + 1, i, preBlackRowCount1))
            return res

        res = dfs(0, 0, 0)
        dfs.cache_clear()
        return res


print(
    Solution().maximumScore(
        grid=[[0, 0, 0, 0, 0], [0, 0, 3, 0, 0], [0, 1, 0, 0, 0], [5, 0, 0, 3, 0], [0, 0, 0, 0, 2]]
    )
)
# grid = [[10,9,0,0,15],[7,1,0,8,0],[5,20,0,11,0],[0,0,0,1,2],[8,12,1,10,3]]
print(
    Solution().maximumScore(
        grid=[
            [10, 9, 0, 0, 15],
            [7, 1, 0, 8, 0],
            [5, 20, 0, 11, 0],
            [0, 0, 0, 1, 2],
            [8, 12, 1, 10, 3],
        ]
    )
)
