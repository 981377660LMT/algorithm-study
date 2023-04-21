# 闭区间[left,right]内的质数个数
# 1≤L≤R≤1e12
# R−L≤500000
# !在范围里筛素数 闭区间内的质数个数
# 区间内的质数个数(分段筛)
# https://cp-algorithms.com/algebra/sieve-of-eratosthenes.html#find-primes-in-range
# !预先生成所有质数到sqrt(upper),再用这些质数标记区间[lower,upper]内的合数

from math import ceil, sqrt


def countPrime(lower: int, upper: int) -> int:
    """[lower, upper]内的质数个数,1<=lower<=upper<=1e12,upper-lower<=500000"""
    isPrime = [True] * (upper - lower + 1)  # P[i] := i+Lが素数か？
    if lower == 1:
        isPrime[0] = False

    last = int(sqrt(upper))
    for fac in range(2, last + 1):
        start = fac * max(ceil(lower / fac), 2) - lower  # !A 以上の最小の fac の倍数
        while start < len(isPrime):
            isPrime[start] = False
            start += fac
    return sum(isPrime)


L, R = map(int, input().split())
print(countPrime(L, R))
