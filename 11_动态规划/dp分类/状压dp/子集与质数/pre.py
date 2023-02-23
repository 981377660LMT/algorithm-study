from collections import Counter
from typing import List


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "_minPrime"  # 每个数的最小质因数

    def __init__(self, maxN: int):
        """预处理 O(nloglogn)"""
        minPrime = list(range(maxN + 1))
        upper = int(maxN**0.5) + 1
        for i in range(2, upper):
            if minPrime[i] < i:
                continue
            for j in range(i * i, maxN + 1, i):
                if minPrime[j] == j:
                    minPrime[j] = i
        self._minPrime = minPrime

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self._minPrime[n] == n

    def getPrimeFactors(self, n: int) -> Counter[int]:
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


MOD = int(1e9 + 7)
E = EratosthenesSieve(100)
N = 30
P = [x for x in range(2, N + 1) if E.isPrime(x)]  # !N以内的质数 (N=30时，质数有10个)
MP = {p: i for i, p in enumerate(P)}  # !每个质数的编号
F = [E.getPrimeFactors(x) for x in range(N + 1)]  # !N以内的数的质因数分解
FMAX = [max(fs.values(), default=0) for fs in F]  # !每个数的质因数分解中最多的质因数的个数
MASK = [sum(1 << MP[p] for p in f) for f in F]  # !每个数的质因数分解对应的mask
