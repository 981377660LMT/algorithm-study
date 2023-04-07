# 给定一个n个点，n 条边的环，有k种颜色，给每个顶点染色，
# 问有多少种本质不同的染色方案，答案对1e9＋7取模
# 注意本题的本质不同，定义为:只需要不能通过旋转与别的染色方案相同。
# n<=1e9
# 不同着色方案的计数问题 (polya定理)


from collections import Counter
from typing import List

MOD = int(1e9 + 7)


def count_cycle_coloring(n: int, k: int) -> int:
    """统计环上的着色方案数"""
    factors = getPrimeFactors(n)
    divs = getFactors(n)
    res = 0
    for div in divs:
        e = div
        for p in factors:
            if div % p == 0:
                e = e // p * (p - 1)
        res += pow(k, n // div, MOD) * e
        res %= MOD
    return res * pow(n, MOD - 2, MOD) % MOD


def getFactors(n: int) -> List[int]:
    """n 的所有因数 O(sqrt(n))"""
    if n <= 0:
        return []
    small, big = [], []
    upper = int(n**0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


def getPrimeFactors(n: int) -> "Counter[int]":
    """n 的素因子分解 O(sqrt(n))"""
    res = Counter()
    upper = int(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i
    if n > 1:
        res[n] += 1
    return res
