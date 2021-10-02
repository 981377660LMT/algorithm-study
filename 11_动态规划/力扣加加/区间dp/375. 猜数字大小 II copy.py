from functools import lru_cache


class Solution:
    def getMoneyAmount(self, n: int) -> int:
        @lru_cache(None)
        def dfs(l, r):
            if l == r:
                return 0
            if r - l == 1:
                return l
            res = float('inf')
            for i in range(l, r + 1):
                res = min(res, i + max(dfs(l, i - 1), dfs(i + 1, r)))
            return res

        return dfs(1, n)

