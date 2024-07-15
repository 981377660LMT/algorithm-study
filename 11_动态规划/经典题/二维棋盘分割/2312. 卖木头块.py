# 2312. 卖木头块
# https://leetcode.cn/problems/selling-pieces-of-wood/description/
# prices[i] = [hi, wi, pricei] 表示你可以以 pricei 元的价格卖一块高为 hi 宽为 wi 的矩形木块
# 所有 (hi, wi) 互不相同 。
# 每一次操作中，你必须按下述方式之一执行切割操作，以得到两块更小的矩形木块：
# 沿垂直方向按高度 完全 切割木块，或
# 沿水平方向按宽度 完全 切割木块
# 请你返回切割一块大小为 m x n 的木块后，能得到的 最多 钱数。
# 1 <= m, n <= 200
# 1 <= prices.length <= 2 * 1e4

from collections import defaultdict
from functools import lru_cache
from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def sellingWood(self, m: int, n: int, prices: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(row: int, col: int) -> int:
            """r行c列的矩形木块可以切出来的最大价值.
            时间复杂度: O(m*n*m+m*n*n)
            """
            res = mapping[(row, col)]
            # 横着切
            for r in range(1, row // 2 + 1):
                res = max2(res, dfs(r, col) + dfs(row - r, col))
            # 竖着切
            for c in range(1, col // 2 + 1):
                res = max2(res, dfs(row, c) + dfs(row, col - c))
            return res

        mapping = defaultdict(int, {(h, w): price for h, w, price in prices})
        res = dfs(m, n)
        dfs.cache_clear()
        return res
