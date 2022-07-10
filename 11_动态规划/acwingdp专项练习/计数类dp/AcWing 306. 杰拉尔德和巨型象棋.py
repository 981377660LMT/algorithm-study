# row,col<=1e5
# 需要用O(nlogn)的解法

# 左上角移动到右下角，一共有多少种路线
# 每一步可以向右或向下移动一格，不可以经过障碍物
# https://www.acwing.com/activity/content/code/content/2328319/
# 容斥原理
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
MOD = int(1e9 + 7)


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


@lru_cache(None)
def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


row, col, n = map(int, input().split())
bad = []
bad.append((0, 0))
for _ in range(n):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    bad.append((x, y))
bad.append((row - 1, col - 1))
bad.sort()
# 原点出发走到第i个障碍物且中间不经过其余障碍物的方案数
# 到当前第i个障碍的合法方案等于到i的所有方案减去到i之前每个障碍的不合法方案
# 终点看作是最后一个障碍

dp = [0] * (n + 10)
dp[0] = 1
for i in range(1, n + 2):
    r, c = bad[i]
    dp[i] = C(r + c - 2, r - 1)
    for j in range(1, i):
        preR, preC = bad[j]
        if preR <= r and preC <= c:
            dp[i] -= dp[j] * C(r + c - preR - preC, r - preR)
            dp[i] %= MOD
print(dp[n + 1])
