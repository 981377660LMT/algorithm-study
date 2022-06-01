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

            return dfs(index + 1, remain) + dfs(index, remain - coins[index])

        n = len(coins)
        res = dfs(0, amount)
        dfs.cache_clear()
        return res
