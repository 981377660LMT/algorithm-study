"""primes"""

from collections import Counter
from functools import lru_cache
from math import floor
from typing import List


def genPrime(n: int) -> List[int]:
    """筛法求小于等于n的素数"""
    isPrime = [True] * (n + 1)
    res = []
    for num in range(2, n + 1):
        if isPrime[num]:
            res.append(num)
            for multi in range(num * num, n + 1, num):
                isPrime[multi] = False
    return res


@lru_cache(None)
def getPrimeFactors(n: int) -> Counter:
    """n 的质因数分解 """
    res = Counter()
    upper = floor(n ** 0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


def getFactors(n: int) -> List[int]:
    """n 的所有因数"""
    if n <= 0:
        return []
    small, big = [], []
    upper = floor(n ** 0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


def isPrime(n: int) -> bool:
    """判断n是否是素数"""
    if n < 2:
        return False
    upper = floor(n ** 0.5) + 1
    for i in range(2, upper):
        if n % i == 0:
            return False
    return True
