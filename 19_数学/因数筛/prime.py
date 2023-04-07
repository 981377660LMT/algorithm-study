"""primes"""

from collections import Counter, defaultdict
from functools import lru_cache
from math import floor, gcd
from random import randint
from typing import DefaultDict, List


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

    def getPrimeFactors(self, n: int) -> "DefaultDict[int, int]":
        """求n的质因数分解 O(logn)"""
        res, f = defaultdict(int), self.minPrime
        while n > 1:
            m = f[n]
            res[m] += 1
            n //= m
        return res

    def getPrimes(self) -> List[int]:
        return [x for i, x in enumerate(self.minPrime) if i >= 2 and i == x]


def getPrimes(n: int) -> List[int]:
    """埃氏筛求小于等于n的素数 O(nloglogn)"""
    isPrime = [True] * (n + 1)
    for num in range(2, n + 1):
        if isPrime[num]:
            for multi in range(num * num, n + 1, num):
                isPrime[multi] = False
    return [num for num in range(2, n + 1) if isPrime[num]]


def getFactors(n: int) -> List[int]:
    """n 的所有因数 O(sqrt(n))"""
    if n <= 0:
        return []
    small, big = [], []
    upper = int(n**0.5) + 1
    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            if i != n // i:
                big.append(n // i)
    return small + big[::-1]


def isPrime(n: int) -> bool:
    """判断n是否是素数 O(sqrt(n))"""
    if n < 2:
        return False
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        if n % i == 0:
            return False
    return True


@lru_cache(None)
def getPrimeFactors1(n: int) -> "Counter[int]":
    """n 的素因子分解 O(sqrt(n))"""
    res = Counter()
    upper = int(n**0.5) + 1  # isqrt(n) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i

    # 注意考虑本身
    if n > 1:
        res[n] += 1
    return res


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


def countPrimes(n: int) -> int:
    """
    计算不超过n的素数个数
    1<=n<=1e11

    质数的数目为 π(n) = O(n/logn)
    """
    if n < 2:
        return 0
    v = int(n**0.5) + 1
    smalls = [i // 2 for i in range(1, v + 1)]
    smalls[1] = 0
    s = v // 2
    roughs = [2 * i + 1 for i in range(s)]
    larges = [(n // (2 * i + 1) + 1) // 2 for i in range(s)]
    skip = [False] * v

    pc = 0
    for p in range(3, v):
        if smalls[p] <= smalls[p - 1]:
            continue

        q = p * p
        pc += 1
        if q * q > n:
            break
        skip[p] = True
        for i in range(q, v, 2 * p):
            skip[i] = True

        ns = 0
        for k in range(s):
            i = roughs[k]
            if skip[i]:
                continue
            d = i * p
            larges[ns] = larges[k] - (larges[smalls[d] - pc] if d < v else smalls[n // d]) + pc
            roughs[ns] = i
            ns += 1
        s = ns
        for j in range((v - 1) // p, p - 1, -1):
            c = smalls[j] - pc
            e = min((j + 1) * p, v)
            for i in range(j * p, e):
                smalls[i] -= c

    for k in range(1, s):
        m = n // roughs[k]
        s = larges[k] - (pc + k - 1)
        for l in range(1, k):
            p = roughs[l]
            if p * p > m:
                break
            s -= smalls[m // p] - (pc + l - 1)
        larges[0] -= s

    return larges[0]


if __name__ == "__main__":
    for i in range(1000):
        assert getPrimeFactors1(i) == getPrimeFactors2(i)

    MOD = int(1e9 + 7)
    fac = [1, 1, 2]  # 阶乘打表
    while len(fac) <= 100:
        fac.append(fac[-1] * len(fac) % MOD)

    class Solution:
        def numPrimeArrangements(self, n: int) -> int:
            def countPrime(upper: int) -> int:
                """统计[1, upper]中的素数个数 upper<=1e5"""
                isPrime = [True] * (upper + 1)
                res = 0
                for num in range(2, upper + 1):
                    if isPrime[num]:
                        res += 1
                        for mul in range(num * num, upper + 1, num):
                            isPrime[mul] = False
                return res

            ok = countPrime(n)
            ng = n - ok
            return (fac[ok] * fac[ng]) % MOD
