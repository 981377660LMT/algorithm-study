from functools import lru_cache
import gc
from typing import List

gc.disable()

# 1 <= matchsticks.length <= 15
# 1 <= matchsticks[i] <= 108


class Solution:
    def makesquare(self, matchsticks: List[int]) -> bool:
        @lru_cache(None)
        def dfs(visited: int, curSum: int) -> bool:
            """注意这里curSum由visited唯一确定 因此复杂度是O(n*2^n)"""
            if visited == target:
                return True
            for i in range(n):
                if (visited >> i) & 1:
                    continue
                if curSum + matchsticks[i] <= quarter:
                    # 达到quarter表示可以组成一条边
                    if dfs(visited | (1 << i), (curSum + matchsticks[i]) % quarter):
                        return True
            return False

        n = len(matchsticks)
        if n < 4:
            return False

        sum_ = sum(matchsticks)
        if sum_ % 4 != 0:
            return False
        quarter = sum_ // 4

        target = (1 << n) - 1
        res = dfs(0, 0)
        dfs.cache_clear()
        return res

