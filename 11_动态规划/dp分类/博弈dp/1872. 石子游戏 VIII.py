from functools import lru_cache
from itertools import accumulate
from typing import List

# 2 <= n <= 1e5


class Solution:
    def stoneGameVIII(self, stones: List[int]) -> int:
        """to jump or not to jump
        
        线性dp
        """

        @lru_cache(None)
        def dfs(index: int) -> int:
            """选择前index个石头,最大的收益 每次至少移除两个石子"""
            if index >= n:
                return preSum[n]
            return max(dfs(index + 1), preSum[index] - dfs(index + 1))  # to jump or not to jump

        n = len(stones)
        preSum = list(accumulate(stones, initial=0))
        res = dfs(2)
        dfs.cache_clear()
        return res

