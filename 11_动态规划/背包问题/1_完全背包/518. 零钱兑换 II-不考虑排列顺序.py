# 请你计算并返回可以凑成总金额的硬币组合数。如果任何硬币组合都无法凑出总金额，返回 0 。
# 假设每一种面额的硬币有无限个。

from functools import lru_cache
from typing import List


class Solution:
    def change(self, amount: int, coins: List[int]) -> int:
        """不考虑排列顺序 每个物品选多少个(jump or not)"""

        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain <= 0:
                return int(remain == 0)
            if index == n:
                return int(remain == 0)

            # !to jump or not to jump
            return dfs(index + 1, remain) + dfs(index, remain - coins[index])

        n = len(coins)
        res = dfs(0, amount)
        dfs.cache_clear()
        return res
