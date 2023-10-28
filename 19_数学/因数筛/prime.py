"""primes"""

from typing import DefaultDict, List, Mapping, Tuple
from collections import Counter, defaultdict
from math import ceil, floor, gcd, sqrt
from random import randint


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


def countFactors(primeFactors: "Mapping[int, int]") -> int:
    """
    返回约数个数.`primeFactors`为这个数的所有质数因子分解.
    如果`primeFactors`为空,返回1.
    """
    res = 1
    for count in primeFactors.values():
        res *= count + 1
    return res


def countFactors2(n: int) -> int:
    if n <= 0:
        return 0
    res = 1
    if n & 1 == 0:
        e = 2
        n >>= 1
        while n & 1 == 0:
            n >>= 1
            e += 1
        res *= e
    f = 3
    while f * f <= n:
        if n % f == 0:
            e = 2
            n //= f
            while n % f == 0:
                n //= f
                e += 1
            res *= e
        f += 2
    if n > 1:
        res *= 2
    return res


def countFactorsOfAll(upper: int) -> List[int]:
    """返回[0,upper]的所有数的约数个数."""
    res = [0] * (upper + 1)
    for i in range(1, upper + 1):
        for j in range(i, upper + 1, i):
            res[j] += 1
    return res


def sumFactors(primeFactors: "Mapping[int, int]") -> int:
    """
    返回约数之和.`primeFactors`为这个数的所有质数因子分解.
    如果`primeFactors`为空,返回1.
    """
    res = 1
    for p, count in primeFactors.items():
        cur = 1
        for _ in range(count):
            cur = cur * p + 1
        res *= cur
    return res


def sumFactors2(n: int) -> int:
    if n <= 0:
        return 0
    res = 1
    if n & 1 == 0:
        cur = 1
        while n & 1 == 0:
            n >>= 1
            cur = cur * 2 + 1
        res *= cur
    f = 3
    while f * f <= n:
        if n % f == 0:
            cur = 1
            while n % f == 0:
                n //= f
                cur = cur * f + 1
            res *= cur
        f += 2
    if n > 1:
        res *= n + 1
    return res


def sumFactorsOfAll(upper: int) -> List[int]:
    """返回[0,upper]的所有数的约数之和."""
    res = [0] * (upper + 1)
    for i in range(1, upper + 1):
        for j in range(i, upper + 1, i):
            res[j] += i
    return res


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


def getFactorsOfAll(upper: int) -> List[List[int]]:
    """返回区间 `[0, upper]` 内所有数的约数."""
    res = [[] for _ in range(upper + 1)]
    for i in range(1, upper + 1):
        for j in range(i, upper + 1, i):
            res[j].append(i)
    return res


def isPrime(n: int) -> bool:
    """判断n是否是素数 O(sqrt(n))"""
    if n < 2 or (n >= 4 and n % 2 == 0):
        return False
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        if n % i == 0:
            return False
    return True


# https://judge.yosupo.jp/problem/primality_test
def isPrimeFast(n: int) -> bool:
    "O(logN) miller rabin algorithm"
    if n == 2:
        return True
    if n == 1 or not n & 1:
        return False
    # miller_rabin
    if n < 1 << 30:
        tests = [2, 7, 61]
    else:
        tests = [2, 325, 9375, 28178, 450775, 9780504, 1795265022]
    d = n - 1
    while ~d & 1:
        d >>= 1
    for a in tests:
        if n <= a:
            break
        t = d
        y = pow(a, t, n)
        while t != n - 1 and y != 1 and y != n - 1:
            y = y * y % n
            t <<= 1
        if y != n - 1 and not t & 1:
            return False
    return True


def getPrimeFactors(n: int) -> "Counter[int]":
    res = Counter()
    if n <= 1:
        return res

    count2 = 0
    while n & 1 == 0:
        n >>= 1
        count2 += 1
    if count2:
        res[2] = count2

    cur = 3
    while cur * cur <= n:
        count = 0
        while n % cur == 0:
            n //= cur
            count += 1
        if count:
            res[cur] = count
        cur += 2

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


def countPrime(lower: int, upper: int) -> int:
    """[lower, upper]内的质数个数,1<=lower<=upper<=1e12,upper-lower<=500000"""
    isPrime = [True] * (upper - lower + 1)  # P[i] := i+L是否为质数
    if lower == 1:
        isPrime[0] = False

    last = int(sqrt(upper))
    for fac in range(2, last + 1):
        start = fac * max(ceil(lower / fac), 2) - lower  # !>=lower的最小fac的倍数
        while start < len(isPrime):
            isPrime[start] = False
            start += fac
    return sum(isPrime)


# 区间筛/区间素数
def segmentedSieve(floor: int, higher: int) -> List[bool]:
    """分段筛求 [floor,higher) 中的每个数是否为质数.
    1<=floor<=higher<=1e12,higher-floor<=5e5
    """
    root = 1
    while (root + 1) * (root + 1) < higher:
        root += 1
    is_prime = [True] * (root + 1)
    is_prime[0] = False
    is_prime[1] = False
    res = [True] * (higher - floor)
    if floor < 2:
        res[: 2 - floor] = [False] * (2 - floor)
    for i in range(2, root + 1):
        if is_prime[i]:
            for j in range(i * i, root + 1, i):
                is_prime[j] = False
            for j in range(max((floor + i - 1) // i, 2) * i, higher, i):
                res[j - floor] = False
    return res


def maxDivisorNum(n: int) -> Tuple[int, int]:
    """n 以内的最多约数个数，以及对应的最小数字.
    n <= 1e9
    """

    primes = (
        2,
        3,
        5,
        7,
        11,
        13,
        17,
        19,
        23,
        29,
        31,
        37,
        41,
        43,
        47,
        53,
        59,
        61,
        67,
        71,
        73,
        79,
        83,
        89,
        97,
    )
    count, res = 0, 1

    def dfs(i: int, maxExp: int, curCount: int, curRes: int) -> None:
        nonlocal count, res
        if (curCount > count) or (curCount == count and curRes < res):
            count, res = curCount, curRes
        for e in range(1, maxExp + 1):
            curRes *= primes[i]
            if curRes > n:
                break
            dfs(i + 1, e, curCount * (e + 1), curRes)

    dfs(0, n.bit_length(), 1, 1)
    return count, res


def maxDivisorNumWithLimit(maxCount: int) -> Tuple[int, int]:
    """在有 最大约数个数限制 的前提下, maxCount 最大是多少，以及对应的最小数字."""
    if maxCount == 0:
        return 0, 0
    left, right = 0, int(1e9)
    while left <= right:
        mid = (left + right) // 2
        count, _ = maxDivisorNum(mid)
        if count > maxCount:
            right = mid - 1
        else:
            left = mid + 1
    return maxDivisorNum(right)


def maxDivisorNumInInterval(min: int, max: int) -> Tuple[int, int]:
    """[min,max]以内的最多约数个数，以及对应的最小数字.
    1<=min<=max<=1e9
    """
    if max - min <= 100000:
        count, res = 0, 0
        for i in range(min, max + 1):
            curCount = countFactors2(i)
            if curCount > count:
                count, res = curCount, i
        return count, res

    primes = (
        2,
        3,
        5,
        7,
        11,
        13,
        17,
        19,
        23,
        29,
        31,
        37,
        41,
        43,
        47,
        53,
        59,
        61,
        67,
        71,
        73,
        79,
        83,
        89,
        97,
    )
    count, res = 0, 0

    def dfs(i: int, maxExp: int, curCount: int, curRes: int) -> None:
        nonlocal count, res
        if curRes >= min and (curCount > count or (curCount == count and curRes < res)):
            count, res = curCount, curRes
        for e in range(1, maxExp + 1):
            curRes *= primes[i]
            if curRes > max:
                break
            dfs(i + 1, e, curCount * (e + 1), curRes)

    dfs(0, max.bit_length(), 1, 1)
    return count, res


def oddDivisorsNum(n: int) -> int:
    """n 拆分成若干连续整数的方法数/奇约数个数"""
    res = 0
    upper = int(sqrt(n)) + 1
    for i in range(1, upper):
        if n % i == 0:
            if i & 1 == 1:
                res += 1
            if i * i < n and n // i & 1 == 1:
                res += 1
    return res


def medianDivisor(n: int) -> int:
    """因子的中位数（偶数个因子时取小的那个）"""
    start = int(sqrt(n))
    for d in range(start, 0, -1):
        if n % d == 0:
            return d
    raise ValueError("medianDivisor: n must be positive")


if __name__ == "__main__":
    for i in range(1, 1000):
        assert getPrimeFactors(i) == getPrimeFactors2(i)
        assert countFactors2(i) == len(getFactors(i)) == countFactors(getPrimeFactors(i))
        assert sumFactors(getPrimeFactors(i)) == sumFactors2(i) == sum(getFactors(i))
    for i in range(10 + 1):
        print(maxDivisorNumWithLimit(i))
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

    import time

    time1 = time.time()
    getPrimeFactors(int(1e14))
    print(time.time() - time1)
