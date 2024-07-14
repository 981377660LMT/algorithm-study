# 100361. 切蛋糕的最小总开销 I
# https://leetcode.cn/problems/minimum-cost-for-cutting-cake-i/description/
# 有一个 m x n 大小的矩形蛋糕，需要切成 1 x 1 的小块。
# 给你整数 m ，n 和两个数组：
# horizontalCut 的大小为 m - 1 ，其中 horizontalCut[i] 表示沿着水平线 i 切蛋糕的开销。
# verticalCut 的大小为 n - 1 ，其中 verticalCut[j] 表示沿着垂直线 j 切蛋糕的开销。
# 一次操作中，你可以选择任意不是 1 x 1 大小的矩形蛋糕并执行以下操作之一：
# 沿着水平线 i 切开蛋糕，开销为 horizontalCut[i] 。
# 沿着垂直线 j 切开蛋糕，开销为 verticalCut[j] 。
# 每次操作后，这块蛋糕都被切成两个独立的小蛋糕。
# 每次操作的开销都为最开始对应切割线的开销，并且不会改变。
# 请你返回将蛋糕全部切成 1 x 1 的蛋糕块的 最小 总开销。
# 1 <= m, n <= 20
# horizontalCut.length == m - 1
# verticalCut.length == n - 1
# 1 <= horizontalCut[i], verticalCut[i] <= 1e3
#
# O(m^2*n^2*(m+n))

from functools import lru_cache
from typing import List

INF = int(1e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def minimumCost(self, m: int, n: int, horizontalCut: List[int], verticalCut: List[int]) -> int:
        @lru_cache(None)
        def dfs(row1: int, col1: int, row2: int, col2: int) -> int:
            if row1 == row2 and col1 == col2:
                return 0
            res = INF
            # 横着切
            for i in range(row1, row2):
                cur = dfs(row1, col1, i, col2) + dfs(i + 1, col1, row2, col2) + horizontalCut[i]
                res = min2(res, cur)
            # 竖着切
            for i in range(col1, col2):
                cur = dfs(row1, col1, row2, i) + dfs(row1, i + 1, row2, col2) + verticalCut[i]
                res = min2(res, cur)
            return res

        res = dfs(0, 0, m - 1, n - 1)
        dfs.cache_clear()
        return res
