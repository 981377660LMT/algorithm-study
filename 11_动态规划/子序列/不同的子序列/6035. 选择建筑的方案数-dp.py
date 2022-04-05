from functools import lru_cache

# 常数很大的O(n)解法


class Solution:
    def numberOfWays(self, s: str) -> int:
        @lru_cache(maxsize=None, typed=False)
        def dfs(index: int, pre: str, count: int) -> int:
            if count == 3:
                return 1
            if index >= len(s):
                return 0

            res = dfs(index + 1, pre, count)
            if (pre == '0' and s[index] == '1') or (pre == '1' and s[index] == '0'):
                res += dfs(index + 1, s[index], count + 1)

            return res

        res = 0
        for start in range(len(s)):
            res += dfs(start, s[start], 1)
        dfs.cache_clear()
        return res


print(Solution().numberOfWays(s="001101"))
