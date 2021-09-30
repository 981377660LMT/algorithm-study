import functools


class Solution:
    def countArrangement(self, N: int) -> int:
        @functools.lru_cache(None)
        def dfs(i, visited):
            if i == N + 1:
                return 1
            res = 0
            for k in range(1, N + 1):
                num = 1 << k
                if not visited & num and ((k) % i == 0 or i % (k) == 0):
                    res += dfs(i + 1, visited | num)
            return res

        return dfs(1, 0)

