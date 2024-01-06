from math import isqrt
from typing import Generator, List, Optional, Tuple
from bisect import bisect_left


def min2(a: int, b: int) -> int:
    return a if a <= b else b


class PrimeTable:
    __slots__ = "_primes", "_sieve", "_done"

    _S = 32768

    def __init__(self, initLimit: int):
        initLimit = max(2, initLimit)
        self._primes = [2]
        self._sieve = [False] * (self._S + 1)
        self._done = initLimit
        self._expand(initLimit + 1)

    def getPrimes(self, limit: Optional[int] = None) -> List[int]:
        """
        返回`小于等于`limit的所有素数.
        每次都会返回一个新的list.
        """
        if limit is None:
            limit = self._done
        if limit > self._done:
            self._expand(limit + 1)
        pos = bisect_left(self._primes, limit + 1)
        return self._primes[:pos]

    def enumerateRangePrimeFactors(
        self, start: int, end: int
    ) -> Generator[Tuple[int, int], None, None]:
        """
        遍历区间[start, end)内所有数的所有素因子.
        生成器返回值: (n, factor), 其中n是[start, end)内的数, factor是n的一个素因子.
        """
        n = end - start
        primes = self._primes
        res = [start + i for i in range(n)]
        pos = bisect_left(primes, isqrt(end) + 1)
        for i in range(pos):
            p = primes[i]
            pp = 1
            while pp <= end // p:
                pp *= p
                s = ((start + pp - 1) // pp) * pp
                while s < end:
                    yield s, p
                    res[s - start] //= p
                    s += pp
        for i, v in enumerate(res):
            if v > 1:
                yield start + i, v

    def enumeratePrefixPrimeFactors(self, n: int) -> Generator[Tuple[int, int], None, None]:
        """
        遍历[2, n]内所有数的所有素因子.
        生成器返回值: (n, factor), 其中n是[2, n]内的数, factor是n的一个素因子.
        """
        pos = bisect_left(self._primes, n + 1)
        primes = self._primes
        for i in range(pos):
            p = primes[i]
            for x in range(n // p, 0, -1):
                yield x * p, p

    def _expand(self, limit: int) -> None:
        if self._done >= limit:
            return
        R = limit // 2
        self._sieve = [False] * (self._S + 1)
        self._primes = [2]
        cp = []
        for i in range(3, self._S + 1, 2):
            if not self._sieve[i]:
                cp.append([i, i * i // 2])
                for j in range(i * i, self._S + 1, 2 * i):
                    self._sieve[j] = True
        for L in range(1, R + 1, self._S):
            block = [False] * self._S
            for i in range(len(cp)):
                p, idx = cp[i]
                ptr = idx
                while ptr < self._S + L:
                    block[ptr - L] = True
                    ptr += p
                cp[i][1] = ptr
            for i in range(min2(self._S, R - L)):
                if not block[i]:
                    self._primes.append((L + i) * 2 + 1)


if __name__ == "__main__":
    # https://leetcode.cn/problems/count-primes/
    class Solution:
        def countPrimes(self, n: int) -> int:
            return len(PrimeTable(n - 1).getPrimes())

    # https://atcoder.jp/contests/abc227/tasks/abc227_g
    # 求C(n,k)的正约数个数, 模998244353.
    # n<=1e12,k<=1e6.
    def abc227g():
        from sys import stdin

        input = stdin.readline

        MOD = 998244353

        table = PrimeTable(int(1e6 + 10))

        n, k = map(int, input().split())
        counter = dict()
        for _, factor in table.enumerateRangePrimeFactors(1, k + 1):
            counter[factor] = counter.get(factor, 0) - 1
        for _, factor in table.enumerateRangePrimeFactors(n - k + 1, n + 1):
            counter[factor] = counter.get(factor, 0) + 1
        res = 1
        for v in counter.values():
            res = (res * (v + 1)) % MOD
        print(res)

    abc227g()
