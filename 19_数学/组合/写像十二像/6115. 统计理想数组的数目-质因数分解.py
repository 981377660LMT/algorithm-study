"""
https://atcoder.jp/contests/arc116/tasks/arc116_c
n,ai<=1e4
序列每个数都是后一个数的因子
求序列个数

长度为n结尾为x的情况怎么计算呢
对x的每个质因子计算 质因子p的个数为k
n-1个物品放到k+1个槽(1,p,p^2,...p^k)里 槽可空 

分解质因数+n个物品放到k个槽的方案数
"""

from collections import Counter
from functools import lru_cache
from math import floor

MOD = int(1e9 + 7)


class Solution:
    def idealArrays(self, n: int, maxValue: int) -> int:
        res = 0
        for end in range(1, maxValue + 1):
            cur = 1
            counter = getPrimeFactors(end)
            for _, k in counter.items():
                # !n-1个物品放到k+1个槽(1,p,p^2,...p^k)里 即选k个隔板 隔板可重复选
                # cur *= CWithReplacement(n, k)  # 和下面一样的
                cur *= put(n - 1, k + 1)
                cur %= MOD
            res += cur
            res %= MOD
        return res


@lru_cache(None)
def getPrimeFactors(
    n: int,
) -> Counter[int]:
    """返回 n 的质因子分解"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    if n > 1:
        res[n] += 1
    return res


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
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


def CWithReplacement(n: int, k: int) -> int:
    """可重复选取的组合数 itertools.combinations_with_replacement 的个数"""
    return C(n + k - 1, k)


def put(n: int, k: int) -> int:
    """
    n个物品放入k个槽(槽可空)的方案数
    """
    return C(n + k - 1, k - 1)


print(Solution().idealArrays(n=2, maxValue=5))  # 10
print(Solution().idealArrays(n=5, maxValue=3))  # 11
print(Solution().idealArrays(n=5, maxValue=9))  # 111
