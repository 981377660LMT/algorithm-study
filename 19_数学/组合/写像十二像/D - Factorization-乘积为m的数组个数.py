# D - Factorization-乘积为m的长为n的数组个数
# !n<=1e5 m<=1e9
# 分解质因数 然后就是质数小球放箱子的问题

from collections import Counter
from math import floor


MOD = int(1e9 + 7)


def factorization(n: int, m: int) -> int:
    """A[0]*A[1]*...*A[n-1] = m, A[i] >= 1, 求A的个数模MOD"""
    primes = getPrimeFactors1(m)
    res = 1
    for v in primes.values():
        res *= put(v, n)
        res %= MOD
    return res


def getPrimeFactors1(n: int) -> "Counter[int]":
    """n 的素因子分解 O(sqrt(n))"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


N = int(2e5 + 10)
fac = [1] * N
ifac = [1] * N
for i in range(1, N):
    fac[i] = (fac[i - 1] * i) % MOD
    ifac[i] = (ifac[i - 1] * pow(i, MOD - 2, MOD)) % MOD


def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac[n] * ifac[k]) % MOD * ifac[n - k]) % MOD


def put(n: int, k: int) -> int:
    """n个物品放入k个槽(槽可空)的方案数"""
    return C(n + k - 1, k - 1)


n, m = map(int, input().split())
print(factorization(n, m))
