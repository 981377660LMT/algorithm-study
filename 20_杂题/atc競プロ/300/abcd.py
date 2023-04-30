from bisect import bisect_right
from math import sqrt
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 以下の正整数のうち、
# a<b<c なる 素数
# a,b,c を用いて
# a
# 2
#  ×b×c
# 2
#   と表せるものはいくつありますか?


# N は
# 300≤N≤10
# 12
# を満たす整数
# a最大是多少


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


S = EratosthenesSieve(int(1e6 + 10))

P = S.getPrimes()
if __name__ == "__main__":
    N = int(input())
    res = 0
    for i, a in enumerate(P):
        for j in range(i + 1, len(P)):
            b = P[j]
            mul = a * a * b
            upper = int(sqrt(N // mul))
            if upper <= b:
                break
            # (b,upper]内的素数个数
            res += bisect_right(P, upper) - j - 1
    print(res)
