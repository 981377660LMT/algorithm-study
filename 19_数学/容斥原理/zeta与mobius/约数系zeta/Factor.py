# 质因数分解

from collections import Counter
from math import gcd
from random import randint


class Factor:
    @staticmethod
    def getPrimeFactors(n: int) -> "Counter[int]":
        """n 的质因数分解 基于PR算法 O(n^1/4*logn)"""
        res = Counter()
        while n > 1:
            p = Factor._PollardRho(n)
            while n % p == 0:
                res[p] += 1
                n //= p
        return res

    @staticmethod
    def _MillerRabin(n: int, k: int = 10) -> bool:
        """米勒-拉宾素性检验(MR)算法判断n是否是素数 O(k*logn*logn)"""
        if n == 2 or n == 3:
            return True
        if n < 2 or n % 2 == 0:
            return False
        d, s = n - 1, 0
        while d % 2 == 0:
            d //= 2
            s += 1
        for _ in range(k):
            a = randint(2, n - 2)
            x = pow(a, d, n)
            if x == 1 or x == n - 1:
                continue
            for _ in range(s - 1):
                x = pow(x, 2, n)
                if x == n - 1:
                    break
            else:
                return False
        return True

    @staticmethod
    def _PollardRho(n: int) -> int:
        """_PollardRho(PR)算法求n的一个因数 O(n^1/4)"""
        if n % 2 == 0:
            return 2
        if n % 3 == 0:
            return 3
        if Factor._MillerRabin(n):
            return n

        x, c = randint(1, n - 1), randint(1, n - 1)
        y, res = x, 1
        while res == 1:
            x = (x * x % n + c) % n
            y = (y * y % n + c) % n
            y = (y * y % n + c) % n
            res = gcd(abs(x - y), n)

        return res if Factor._MillerRabin(res) else Factor._PollardRho(n)  # !这里规定要返回一个素数


if __name__ == "__main__":
    f = 31415926535
    print(Factor.getPrimeFactors(f))
