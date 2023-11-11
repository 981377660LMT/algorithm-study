# 正整数拆分 拆分正整数
from functools import lru_cache
import sys


MOD = int(1e9 + 7)


def splitNum(n: int) -> int:
    """正整数拆分的方案数 完全背包"""
    dp = [0] * (n + 1)
    dp[0] = 1
    for num in range(1, n + 1):
        for cap in range(1, n + 1):
            if cap >= num:
                dp[cap] += dp[cap - num]
                dp[cap] %= MOD
    return dp[n]


print(splitNum(2000))


sys.setrecursionlimit(int(1e6))


def splitNum2(n: int) -> int:
    """正整数拆分的方案数"""

    @lru_cache(None)
    def dfs(index: int, remain: int) -> int:
        if remain < 0:
            return 0
        if index == n + 1:
            return 1 if remain == 0 else 0

        res = 0
        for select in range(remain + 1):
            sum_ = select * index
            if sum_ > remain:
                break
            res += dfs(index + 1, remain - sum_)
            res %= MOD
        return res

    return dfs(1, n)


print(splitNum2(50))  # 204226
