from collections import Counter
from functools import lru_cache
from math import isqrt
from typing import List

# !所有因数
def getFactors(n: int) -> List[int]:
    """返回 n 的所有因数"""
    upper = isqrt(n) + 1
    small, big = [], []

    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            big.append(n // i)

    if small[-1] == big[-1]:
        small.pop()
    return small + big[::-1]


# !所有质因数


@lru_cache(None)
def getPrimeFactors(n: int) -> Counter:
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = isqrt(n) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res
