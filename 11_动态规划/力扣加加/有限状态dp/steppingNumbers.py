# steppingNumber是指相邻位差1的数
# 求n位steppingNumber的个数
from functools import lru_cache


class Solution:
    def solve(self, n):
        @lru_cache(None)
        def dfs(cur: int, pre: int) -> int:
            if cur == n:
                return 1

            res = 0
            for next in {pre + 1, pre - 1}:
                if 0 <= next <= 9:
                    res += dfs(cur + 1, next)
                    res %= int(1e9 + 7)
            return res

        if n == 1:
            return 10

        res = 0
        for start in range(1, 10):
            res += dfs(1, start)
            res %= int(1e9 + 7)
        dfs.cache_clear()
        return res

