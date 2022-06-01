# 这里有 n 个一样的骰子，每个骰子上都有 k 个面，分别标号为 1 到 k 。求掷出target的方案数
from functools import lru_cache


MOD = int(1e9 + 7)


class Solution:
    def numRollsToTarget(self, n: int, k: int, target: int) -> int:
        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            if remain < 0:
                return 0
            if index == n:
                return 1 if remain == 0 else 0

            res = 0
            for i in range(1, k + 1):
                res += dfs(index + 1, remain - i)
                res %= MOD
            return res

        res = dfs(0, target)
        dfs.cache_clear()
        return res

