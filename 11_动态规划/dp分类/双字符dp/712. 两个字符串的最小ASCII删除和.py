from functools import lru_cache


class Solution:
    def minimumDeleteSum(self, s1: str, s2: str) -> int:
        @lru_cache(None)
        def dfs(i: int, j: int) -> int:
            if i == len(s1) and j == len(s2):
                return 0

            res = int(1e20)
            if i < len(s1) and j < len(s2) and s1[i] == s2[j]:
                res = min(res, dfs(i + 1, j + 1))
            if i < len(s1):
                res = min(res, dfs(i + 1, j) + ord(s1[i]))
            if j < len(s2):
                res = min(res, dfs(i, j + 1) + ord(s2[j]))

            return res

        return dfs(0, 0)

