from functools import lru_cache

INF = 0x3FFFFFFF


class Solution:
    def getMoneyAmount(self, n: int) -> int:
        @lru_cache(None)
        def dfs(left, right) -> int:
            if left == right:
                return 0
            if right - left == 1:
                return left

            res = INF
            for i in range(left, right + 1):
                res = min(res, i + max(dfs(left, i - 1), dfs(i + 1, right)))
            return res

        return dfs(1, n)


print(Solution().getMoneyAmount(n=10))
