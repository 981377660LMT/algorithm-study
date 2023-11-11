# 求符号题意的整数个数
# 1. 各位数字之和为k
# 2. 9的倍数
# 3. 各位都是正整数

# k<=1e5
# !9的倍数的性质:各位数字之和为k的9的倍数
# !则有dp[i]=dp[i-1]+dp[i-2]+...+dp[i-9] 打表即可

from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)

k = int(input())
if k % 9 != 0:
    print(0)
    exit(0)


@lru_cache(None)
def dfs(remain: int) -> int:
    if remain <= 0:
        return int(remain == 0)

    res = 0
    for cur in range(1, 10):
        res += dfs(remain - cur)
        res %= MOD
    return res


print(dfs(k))

