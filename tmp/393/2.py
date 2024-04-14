from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


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

    def getPrimeFactors(self, n: int) -> "Counter[int]":
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self._minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self._minPrime) if i >= 2 and i == x]


P = EratosthenesSieve(int(3e5 + 10))

# 给你一个整数数组 nums。


# 返回两个（不一定不同的）素数在 nums 中 下标 的 最大距离。
class Solution:
    def maximumPrimeDifference(self, nums: List[int]) -> int:
        isPrimes = [P.isPrime(x) for x in nums]
        firstPrimeIndex = -1
        lastPrimeIndex = -1
        for i, isPrime in enumerate(isPrimes):
            if isPrime:
                if firstPrimeIndex == -1:
                    firstPrimeIndex = i
                lastPrimeIndex = i
        return lastPrimeIndex - firstPrimeIndex
