from collections import defaultdict
from functools import lru_cache
from typing import List

# prices[i] = [hi, wi, pricei] 表示你可以以 pricei 元的价格卖一块高为 hi 宽为 wi 的矩形木块
# 所有 (hi, wi) 互不相同 。
# 每一次操作中，你必须按下述方式之一执行切割操作，以得到两块更小的矩形木块：
# 沿垂直方向按高度 完全 切割木块，或
# 沿水平方向按宽度 完全 切割木块

# 请你返回切割一块大小为 m x n 的木块后，能得到的 最多 钱数。

# 1 <= m, n <= 200
# 1 <= prices.length <= 2 * 1e4


class Solution:
    def sellingWood(self, m: int, n: int, prices: List[List[int]]) -> int:
        @lru_cache(None)
        def dfs(r: int, c: int) -> int:
            """r行c列的矩形木块可以切出来的最大价值
            
            时间复杂度: O(m*n*m+m*n*n)
            """
            res = mapping[(r, c)]
            for i in range(1, r):
                res = max(res, dfs(i, c) + dfs(r - i, c))  # 垂直切割
            for j in range(1, c):
                res = max(res, dfs(r, j) + dfs(r, c - j))  # 水平切割
            return res

        mapping = defaultdict(int, {(h, w): price for h, w, price in prices})
        res = dfs(m, n)
        dfs.cache_clear()
        return res

