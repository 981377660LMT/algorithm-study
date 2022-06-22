from functools import lru_cache


MOD = int(1e9 + 7)


# 求子序列个数


def numDistinct(s: str, t: str) -> int:
    """求s中有多少个子序列为t，时间复杂度O(st)"""

    @lru_cache(None)
    def dfs(i: int, j: int) -> int:
        if j == len(t):
            return 1
        if i == len(s):
            return 0

        # 选不选当前位置配对 jump or not
        res = dfs(i + 1, j)
        if s[i] == t[j]:
            res += dfs(i + 1, j + 1)
        return res % MOD

    if not t:
        return 0
    return dfs(0, 0) % MOD

