"""primes"""

from collections import Counter
from functools import lru_cache
from math import floor
from typing import List


def getPrimes(n: int) -> List[int]:
    """埃氏筛求小于等于n的素数 O(nloglogn)"""
    isPrime = [True] * (n + 1)
    for num in range(2, n + 1):
        if isPrime[num]:
            for multi in range(num * num, n + 1, num):
                isPrime[multi] = False
    return [num for num in range(2, n + 1) if isPrime[num]]


@lru_cache(None)
def getPrimeFactors(n: int) -> "Counter[int]":
    """n 的质因数分解"""
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


def getFactors(n: int) -> List[int]:
    """n 的所有因数"""
    if n <= 0:
        return []
    small, big = [], []
    upper = floor(n**0.5) + 1
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
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        if n % i == 0:
            return False
    return True


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "_f"

    def __init__(self, maxN: int):
        f = list(range(maxN + 1))
        upper = int(maxN**0.5) + 1
        for i in range(2, upper):
            if f[i] < i:
                continue
            for j in range(i * i, maxN + 1, i):
                if f[j] == j:
                    f[j] = i
        self._f = f

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self._f[n] == n

    def getFactors(self, n: int) -> "Counter[int]":
        """n的质因数分解"""
        res, f = Counter(), self._f
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._f) if i >= 2 and i == x]


if __name__ == "__main__":
    MOD = int(1e9 + 7)
    fac = [1, 1, 2]  # 阶乘打表
    while len(fac) <= 100:
        fac.append(fac[-1] * len(fac) % MOD)

    class Solution:
        def numPrimeArrangements(self, n: int) -> int:
            def countPrime(upper: int) -> int:
                """统计[1, upper]中的素数个数"""
                isPrime = [True] * (upper + 1)
                res = 0
                for num in range(2, upper + 1):
                    if isPrime[num]:
                        res += 1
                        for mul in range(num * num, upper + 1, num):
                            isPrime[mul] = False
                return res

            ok = countPrime(n)
            ng = n - ok
            return (fac[ok] * fac[ng]) % MOD
