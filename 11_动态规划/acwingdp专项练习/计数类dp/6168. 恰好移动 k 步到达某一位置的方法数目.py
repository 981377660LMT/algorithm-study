from functools import lru_cache
from math import comb

MOD = int(1e9 + 7)


class Solution:
    def numberOfWays(self, startPos: int, endPos: int, k: int) -> int:
        # startPos,endPos,k<=1000
        @lru_cache(None)
        def dfs(index: int, remain: int) -> int:
            # !剪枝
            if abs(index - endPos) > remain:
                return 0
            if remain == 0:
                return 1 if index == endPos else 0
            res = dfs(index - 1, remain - 1) + dfs(index + 1, remain - 1)
            return res % MOD

        res = dfs(startPos, k)
        dfs.cache_clear()
        return res

    def numberOfWays2(self, startPos: int, endPos: int, k: int) -> int:
        # 求组合数
        # 向左走x步 向右走y步
        diff = abs(startPos - endPos)
        if (k + diff) & 1:
            return 0
        left = (k + diff) // 2
        return comb(k, left) % MOD
