# 1-n号物品里选若干个物品 要求相邻的物品差要>=k
# 对k=1-n求方案数
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
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


def C(n: int, k: int) -> int:
    if k < 0 or k > n:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


# 从n个物品选a个 有(a-1)*(k-1)个位置不能选
# 因此从n个物品选a个方案数为Comb(n-(a-1)*(k-1),a)
# k逐渐增大 总的复杂度为调和级数 nlogn


n = int(input())
for k in range(1, n + 1):
    res = 0
    select = 1
    while True:
        remain = n - (select - 1) * (k - 1)
        if remain < select:
            break
        res += C(remain, select)
        res %= MOD
        select += 1
    print(res)
