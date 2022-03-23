from functools import lru_cache


MOD = int(1e9 + 7)


# 求子序列个数
class Solution:
    def numDistinct(self, s: str, t: str) -> int:
        """求s中有多少个子序列为t，时间复杂度O(st)"""

        @lru_cache(None)
        def dfs(i1: int, i2: int) -> int:
            if i2 == len(t):
                return 1
            if i1 == len(s):
                return 0

            # 选不选当前位置配对
            if s[i1] == t[i2]:
                return (dfs(i1 + 1, i2 + 1) + dfs(i1 + 1, i2)) % MOD
            return dfs(i1 + 1, i2) % MOD

        if not t:
            return 0

        return dfs(0, 0) % MOD
