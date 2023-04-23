#
#  Divisor function
#
#  Description:
#  sigma(n) = sum[n % d == 0] d
#  equivalently,
#  sigma(p^k) = 1 + p + p^2 + ... + p^k
#  with multiplicative.
#
#  Complexity:
#  divisor_sigma(n):     O(sqrt(n)) by trial division.
#  divisor_sigma(lo,hi): O((hi-lo) loglog(hi)) by prime sieve.
#


# 约数之和函数
# 区间[lo,hi)内所有数的约数之和


from math import ceil
from random import randint
from typing import List


def divisorSum(n: int) -> int:
    """O(sqrt(n))求n的约数之和"""
    res, d = 0, 1
    while d * d < n:
        if n % d == 0:
            res += d + n // d
        d += 1
    if d * d == n:
        res += d
    return res


def divisorSumRange(lo: int, hi: int) -> List[int]:
    """O((hi-lo)loglog(hi)) 求区间[lo,hi)内每个数的约数之和"""
    ps = getPrimes(0, int(hi**0.5) + 1)
    tmp, res = list(range(lo, hi)), [1] * (hi - lo)
    for p in ps:
        k = ((lo + (p - 1)) // p) * p
        while k < hi:
            b = 1
            while tmp[k - lo] > 1 and tmp[k - lo] % p == 0:
                tmp[k - lo] //= p
                b = 1 + b * p
            res[k - lo] *= b
            k += p
    for k in range(lo, hi):
        if tmp[k - lo] > 1:
            res[k - lo] *= 1 + tmp[k - lo]
    return res  # res[k - floor] = res[k]


def getPrimes(lo: int, hi: int) -> List[int]:
    """求区间[lo,hi)内的素数"""

    def max(a, b):
        return a if a > b else b

    def min(a, b):
        return a if a < b else b

    M, SQR = 1 << 14, 1 << 16
    composite, smallComposite = [False] * M, [False] * SQR
    sieve = []
    for i in range(3, SQR, 2):
        if not smallComposite[i]:
            k = i * i + 2 * i * max(0, ceil((lo - i * i) // (2 * i)))
            sieve.append([2 * i, k])
            for j in range(i * i, SQR, 2 * i):
                smallComposite[j] = True
    res = []
    if lo <= 2:
        res.append(2)
        lo = 3
    k = lo | 1
    low = lo
    while low < hi:
        high = min(low + M, hi)
        composite = [False] * M
        for z in sieve:
            while z[1] < high:
                composite[z[1] - low] = True
                z[1] += z[0]
        while k < high:
            if not composite[k - low]:
                res.append(k)
            k += 2
        low += M
    return res


if __name__ == "__main__":
    n = int(1e9)
    left = randint(0, n)
    right = left + 100
    assert divisorSumRange(left, right) == [divisorSum(i) for i in range(left, right)]
    from time import time

    time1 = time()
    print(divisorSumRange(n, n + int(1e5)))
    time2 = time()
    print(time2 - time1)
