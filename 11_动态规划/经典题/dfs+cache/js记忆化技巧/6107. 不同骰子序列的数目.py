from functools import lru_cache
from math import gcd
from typing import Tuple

MOD = int(1e9 + 7)
gcd = lru_cache(gcd)


class Solution:
    def distinctSequences(self, n: int) -> int:
        @lru_cache(None)
        def dfs(index: int, pre: Tuple[int, int]) -> int:
            if index == n:
                return 1

            res = 0
            for cur in set(range(1, 7)) - set(pre):
                if pre[-1] == 0 or gcd(pre[-1], cur) == 1:
                    res += dfs(index + 1, pre[1:] + (cur,))
                    res %= MOD
            return res

        res = dfs(0, tuple([0, 0]))
        dfs.cache_clear()
        return res


print(Solution().distinctSequences(n=4))
print(Solution().distinctSequences(n=2))
print(Solution().distinctSequences(n=200))
