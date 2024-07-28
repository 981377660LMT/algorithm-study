from bisect import bisect_left, bisect_right
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你两个 正整数 l 和 r。对于任何数字 x，x 的所有正因数（除了 x 本身）被称为 x 的 真因数。

# 如果一个数字恰好仅有两个 真因数，则称该数字为 特殊数字。例如：


# 数字 4 是 特殊数字，因为它的真因数为 1 和 2。
# 数字 6 不是 特殊数字，因为它的真因数为 1、2 和 3。
# 返回区间 [l, r] 内 不是 特殊数字 的数字数量。

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


E = EratosthenesSieve(int(1e5))
PS = E.getPrimes()
PS2 = [v * v for v in PS]


# 质数的平方?
class Solution:
    def nonSpecialCount(self, l: int, r: int) -> int:
        upper = bisect_left(PS2, r + 1)
        lower = bisect_left(PS2, l)
        return r - l + 1 - (upper - lower)


# l = 4, r = 16

print(Solution().nonSpecialCount(4, 16))  # 9
