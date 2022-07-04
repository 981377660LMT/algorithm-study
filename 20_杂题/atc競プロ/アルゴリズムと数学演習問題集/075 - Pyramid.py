# 下の段から順に、「隣り合った 2 つの数を足した答えを上の段に書く」という操作を行ったとき、一番上に書かれる整数はいくつですか。
# 金字塔有n层 问金字塔顶的数是多少
# n<=2e5
# !每个数的贡献次数=>二项式定理系数

import sys
from functools import lru_cache


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
    if k < 0 or k > n:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())
nums = list(map(int, input().split()))
res = 0
for i in range(n):
    count = C(n - 1, i)  # 二项式定理系数
    res += count * nums[i]
    res %= MOD
print(res)
