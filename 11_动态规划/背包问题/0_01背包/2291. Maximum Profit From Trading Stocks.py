from functools import lru_cache
import gc
from typing import List

gc.disable()

# 1 <= n <= 1000
# 0 <= present[i], future[i] <= 100
# 0 <= budget <= 1000
# 1e6的数据 python 有一点卡常数


class Solution:
    def maximumProfit(self, A: List[int], B: List[int], k: int) -> int:
        """2564 ms"""

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if index == len(goods):
                return 0
            res = dfs(index + 1, remain)
            cost, score = goods[index]
            if remain >= cost:
                res = max(res, dfs(index + 1, remain - cost) + score)
            return res

        goods = [(cur, next - cur) for cur, next in zip(A, B) if next > cur]
        res = dfs(0, k)
        dfs.cache_clear()
        return res

    def maximumProfit2(self, A: List[int], B: List[int], k: int) -> int:
        """852 ms"""
        goods = [(cur, next - cur) for cur, next in zip(A, B) if next > cur]
        dp = [0] * (k + 1)
        for cost, score in goods:
            for i in range(k, cost - 1, -1):
                dp[i] = max(dp[i], dp[i - cost] + score)
        return dp[-1]
