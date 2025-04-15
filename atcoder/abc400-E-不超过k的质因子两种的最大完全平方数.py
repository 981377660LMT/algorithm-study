# abc400-E-不超过k的质因子两种的最大完全平方数
# https://atcoder.jp/contests/abc400/tasks/abc400_e
#
# ## 问题描述
#
# 对于正整数 N，我们将其称为 400 number，当且仅当满足以下两个条件：
# 1. N 的素因数恰好有两种。
# 2. 对于 N 的每个素因数 P，N 能被 P 整除的次数是偶数次。
#    更严格地说，对于每个素因数 P，使得 P^K 是 N 的约数的最大非负整数 K 是偶数。
# 给定 Q 个查询。每个查询给出一个整数 x，请你求出不超过 x 的最大 400 number。保证对于每个查询，x 以下一定存在 400 number。
#
# ---
#
# ## 约束条件
#
# - 1 <= Q <= 2e5
# !- 对于每个查询，36 <= x <= 1e12
# - 输入的所有值均为整数
#
# !枚举所有的完全平方数，判断是否符合条件，然后二分查找。

from typing import List
from bisect import bisect_right
from collections import Counter


class EratosthenesSieve:
    """埃氏筛"""

    __slots__ = "minPrime"  # 每个数的最小质因数

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
        self.minPrime = minPrime

    def isPrime(self, n: int) -> bool:
        if n < 2:
            return False
        return self.minPrime[n] == n

    def getPrimeFactors(self, n: int) -> "Counter[int]":
        """求n的质因数分解 O(logn)"""
        res, f = Counter(), self.minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self.minPrime) if i >= 2 and i == x]


MAX_N = int(1e6) + 10
E = EratosthenesSieve(MAX_N)

if __name__ == "__main__":

    cands = []
    for v in range(1, MAX_N + 1):
        primeCount = len(E.getPrimeFactors(v))
        if primeCount == 2:
            cands.append(v * v)

    def query(x: int) -> int:
        pos = bisect_right(cands, x)
        return cands[pos - 1]

    Q = int(input())
    for _ in range(Q):
        x = int(input())
        res = query(x)
        print(res)
