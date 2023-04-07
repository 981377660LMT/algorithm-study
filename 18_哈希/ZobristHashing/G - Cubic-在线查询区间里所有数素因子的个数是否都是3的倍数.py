# G - Cubic-在线查询区间里所有数素因子的个数是否都是3的倍数
# 类似棋盘哈希,加上状态个数,即 (值,出现次数)对应一个状态
# https://hackmd.io/@tatyam-prime/r1dg9Q389#

from collections import Counter, defaultdict
from random import randint
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


class AllCountKChecker:
    """判断数据结构中每个数出现的次数是否均k的倍数."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_hash", "_counter", "_k")

    def __init__(self, k: int) -> None:
        self._hash = 0
        self._counter = defaultdict(int)
        self._k = k

    def add(self, x: int) -> None:
        count = self._counter[x]
        self._hash ^= self._poolSingleton[(x, count)]
        count += 1
        if count == self._k:
            count = 0
        self._counter[x] = count
        self._hash ^= self._poolSingleton[(x, count)]

    def remove(self, x: int) -> None:
        """删除前需要保证x在集合中存在."""
        count = self._counter[x]
        self._hash ^= self._poolSingleton[(x, count)]
        count -= 1
        if count == -1:
            count = self._k - 1
        self._counter[x] = count
        self._hash ^= self._poolSingleton[(x, count)]

    def query(self) -> bool:
        return self._hash == 0

    def getHash(self) -> int:
        return self._hash


E = EratosthenesSieve(int(1e6) + 5)
if __name__ == "__main__":
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, q = map(int, input().split())
    nums = list(map(int, input().split()))
    CC = AllCountKChecker(3)
    preHash = [0]
    for x in nums:
        ps = E.getPrimeFactors(x)
        for p, c in ps.items():
            for _ in range(c):
                CC.add(p)
        preHash.append(CC.getHash())

    for _ in range(q):
        l, r = map(int, input().split())
        print("Yes" if preHash[l - 1] == preHash[r] else "No")
