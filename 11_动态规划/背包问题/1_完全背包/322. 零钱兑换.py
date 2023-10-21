# 给你一个整数数组 coins ，表示不同面额的硬币；以及一个整数 amount ，表示总金额。
# !计算并返回可以凑成总金额所需的 最少的硬币个数 。如果没有任何一种硬币组合能组成总金额，返回 -1 。
# 你可以认为每种硬币的数量是无限的。

from functools import lru_cache
from typing import List

INF = int(1e18)


class Solution:
    def coinChange(self, coins: List[int], amount: int) -> int:
        dp = [INF] * (amount + 1)
        dp[0] = 0
        for coin in coins:
            for cur in range(coin, amount + 1):
                dp[cur] = min(dp[cur], dp[cur - coin] + 1)
        return dp[amount] if dp[amount] != INF else -1

    def coinChange2(self, coins: List[int], amount: int) -> int:
        @lru_cache(None)
        def dfs(remain: int) -> int:
            if remain < 0:
                return INF
            if remain == 0:
                return 0
            res = INF
            for cur in coins:
                if remain - cur >= 0:
                    res = min(res, dfs(remain - cur) + 1)
            return res

        res = dfs(amount)
        dfs.cache_clear()
        return res if res != INF else -1
