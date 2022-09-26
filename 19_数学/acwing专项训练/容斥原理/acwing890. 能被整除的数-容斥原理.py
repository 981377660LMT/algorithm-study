"""
容斥原理求解方案数
"""
# 给定一个整数 n 和 m 个不同的质数 p1,p2,…,pm。
# 请你求出 1∼n 中能被 p1,p2,…,pm 中的至少一个数整除的整数有多少个。
# m<=10
# !推广:[a,b]间与num互质的数的个数 1<=a<=b<=10^9 1<=num<=10^9

from collections import Counter
from functools import lru_cache
from math import floor
from typing import List


def count(upper: int, primes: List[int]) -> int:
    """[1, upper]中能被primes中的至少一个数整除的数的个数"""
    m = len(primes)
    res = 0
    for state in range(1, (1 << m)):  # 枚举被哪些数整除
        mul = 1
        for i in range(m):
            if state & (1 << i):
                mul *= primes[i]

        # !奇数个元素系数为 1，偶数个元素为 -1
        if state.bit_count() & 1:
            res += upper // mul
        else:
            res -= upper // mul
    return res


@lru_cache(None)
def getPrimeFactors1(n: int) -> "Counter[int]":
    """n 的素因子分解 O(sqrt(n))"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    if n > 1:
        res[n] += 1
    return res


def solve(lower: int, upper: int, num: int) -> int:
    """[lower, upper]间与num互质的数的个数"""
    primes = list(getPrimeFactors1(num).keys())
    mul = count(upper, primes) - count(lower - 1, primes)
    return upper - lower + 1 - mul


if __name__ == "__main__":
    lower, upper, num = map(int, input().split())
    print(solve(lower, upper, num))
