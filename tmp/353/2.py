from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数 n 。如果两个整数 x 和 y 满足下述条件，则认为二者形成一个质数对：

# 1 <= x <= y <= n
# x + y == n
# x 和 y 都是质数
# 请你以二维有序列表的形式返回符合题目要求的所有 [xi, yi] ，列表需要按 xi 的 非递减顺序 排序。如果不存在符合要求的质数对，则返回一个空数组。


# 注意：质数是大于 1 的自然数，并且只有两个因子，即它本身和 1 。

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


E = EratosthenesSieve(int(1e6 + 10))


class Solution:
    def findPrimePairs(self, n: int) -> List[List[int]]:
        res = []
        for x in range(2, n + 1):
            if not E.isPrime(x):
                continue
            y = n - x
            if x <= y <= n:
                if E.isPrime(y):
                    res.append([x, y])
            else:
                break
        return res
