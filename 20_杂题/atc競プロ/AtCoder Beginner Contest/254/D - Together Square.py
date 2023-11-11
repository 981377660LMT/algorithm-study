# 给你一个n ，求有多少对（i,j）满足i* j是一个平方数, i,j满足小于等于n 。
# n<=2e5

from collections import Counter
import sys
import os
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


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


def main() -> None:
    n = int(input())
    ES = EratosthenesSieve(n)
    match = Counter()
    for i in range(1, n + 1):
        counter = ES.getFactors(i)
        mul = 1
        for p, c in counter.items():
            if c & 1:  # 凑成偶数
                mul *= p
        match[mul] += 1

    res = 0
    for c in match.values():
        res += c * c  # 每种互相配对成偶数
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
