"""
每次操作将n分解为两个因数相乘 素数时停止
递归操作 问最少要多少次 才能停止

每次都会*2 因此是ceil(log(质因数个数))
"""

import sys
from collections import Counter
from functools import lru_cache
from math import ceil, floor, log2


sys.setrecursionlimit(int(1e6))
input = sys.stdin.readline
MOD = int(1e9 + 7)


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


n = int(input())
c = getPrimeFactors(n)
sum_ = sum(c.values())
print(ceil(log2(sum_)))
