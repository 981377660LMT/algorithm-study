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


# 给你两个正整数 left 和 right ，请你找到两个整数 num1 和 num2 ，它们满足：

# left <= nums1 < nums2 <= right  。
# nums1 和 nums2 都是 质数 。
# nums2 - nums1 是满足上述条件的质数对中的 最小值 。
# 请你返回正整数数组 ans = [nums1, nums2] 。如果有多个整数对满足上述条件，请你返回 nums1 最小的质数对。如果不存在符合题意的质数对，请你返回 [-1, -1] 。

# 如果一个整数大于 1 ，且只能被 1 和它自己整除，那么它是一个质数。

S = EratosthenesSieve(int(1e6 + 10))

P = S.getPrimes()


class Solution:
    def closestPrimes(self, left: int, right: int) -> List[int]:
        ok = []
        for p in P:
            if left <= p <= right:
                ok.append(p)
        if len(ok) < 2:
            return [-1, -1]
        res = [ok[0], ok[1]]
        for i in range(1, len(ok)):
            if ok[i] - ok[i - 1] < res[1] - res[0]:
                res = [ok[i - 1], ok[i]]
        return res
