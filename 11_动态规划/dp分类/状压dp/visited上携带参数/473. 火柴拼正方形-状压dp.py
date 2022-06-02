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
            if visited == (1 << n) - 1:
                return True
            for i in range(n):
                if (visited >> i) & 1:
                    continue
                if curSum + matchsticks[i] <= div:
                    # 达到div表示可以组成一条边
                    if dfs(visited | (1 << i), (curSum + matchsticks[i]) % div):
                        return True
            return False

        div, mod = divmod(sum(matchsticks), 4)
        if mod != 0:
            return False
        n = len(matchsticks)
        return dfs(0, 0)

