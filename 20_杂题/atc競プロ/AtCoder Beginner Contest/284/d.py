import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 正整数 N が与えられます。N は、2 つの相異なる素数 p,q を用いて N=p
# 2
#  q と表せることがわかっています。

# p,q を求めてください。

# T 個のテストケースが与えられるので、それぞれについて答えを求めてください。

from collections import Counter, defaultdict
from functools import lru_cache
from math import floor, gcd
from random import randint
from typing import DefaultDict, List


def MillerRabin(n: int, k: int = 10) -> bool:
    """米勒-拉宾素性检验(MR)算法判断n是否是素数 O(k*logn*logn)

    https://zhuanlan.zhihu.com/p/267884783
    """
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


def PollardRho(n: int) -> int:
    """PollardRho(PR)算法求n的一个因数 O(n^1/4)

    https://zhuanlan.zhihu.com/p/267884783
    """
    if n % 2 == 0:
        return 2
    if n % 3 == 0:
        return 3
    if MillerRabin(n):
        return n

    x, c = randint(1, n - 1), randint(1, n - 1)
    y, res = x, 1
    while res == 1:
        x = (x * x % n + c) % n
        y = (y * y % n + c) % n
        y = (y * y % n + c) % n
        res = gcd(abs(x - y), n)

    return res if MillerRabin(res) else PollardRho(n)  # !这里规定要返回一个素数


def getPrimeFactors2(n: int) -> "Counter[int]":
    """n 的质因数分解 基于PR算法 O(n^1/4*logn)"""
    res = Counter()
    while n > 1:
        p = PollardRho(n)
        while n % p == 0:
            res[p] += 1
            n //= p
    return res


# 正整数 N が与えられます。N は、2 つの相異なる素数 p,q を用いて N=p
# 2
#  q と表せることがわかっています。

# p,q を求めてください。

# T 個のテストケースが与えられるので、それぞれについて答えを求めてください。
if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        N = int(input())
        factors = getPrimeFactors2(N)
        res = list(factors.keys())
        a, b = res
        if a**2 * b == N:
            print(a, b)
        else:
            print(b, a)
