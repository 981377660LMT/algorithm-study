# E - Alphabet Tiles
# https://atcoder.jp/contests/abc358/tasks/abc358_e
# 给定n种物品，每种物品有ai个.
# !从这n种物品中选出1-k个物品组成一个序列，求方案数模998244353.

from functools import lru_cache
from Enumeration import Enumeration


MOD = 998244353


C = Enumeration(int(1e4), MOD)


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    K = int(input())
    limits = [int(x) for x in input().split()]
    ##########################################################
    # 记忆化解法

    @lru_cache(None)
    def dfs(gid: int, count: int) -> int:
        """前gid种物品,已经选了count个物品,方案数."""
        if count > K:
            return 0
        if gid == len(limits):
            return 1

        res = 0
        for v in range(limits[gid] + 1):
            if count + v > K:
                break
            res += dfs(gid + 1, count + v) * C.C(count + v, v) % MOD
            res %= MOD
        return res

    res = dfs(0, 0) - 1
    dfs.cache_clear()
    print(res % MOD)

    #########################################################

    limits = [int(x) for x in input().split()]
    dp = [[0] * (K + 1) for _ in range(27)]
    dp[0][0] = 1
    for i in range(1, 27):
        for j in range(K + 1):
            for k in range(limits[i - 1] + 1):
                if j - k >= 0:
                    dp[i][j] += dp[i - 1][j - k] * C.C(j, k)
                    dp[i][j] %= MOD
    print((sum(dp[26]) - 1) % MOD)
