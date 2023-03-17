# 约数之和
# https://www.acwing.com/solution/content/16981/

from collections import Counter
from math import floor

MOD = int(1e9 + 7)


def getPrimeFactors(n: int) -> "Counter[int]":
    """返回 n 的所有质数因子"""
    res = Counter()
    upper = floor(n**0.5) + 1
    for i in range(2, upper):
        while n % i == 0:
            res[i] += 1
            n //= i
    if n > 1:
        res[n] += 1
    return res


def sumOfFactors(counter: "Counter[int]") -> int:
    """返回所有约数之和, counter 为这个数的所有质数因子分解."""
    res = 1
    for p, count in counter.items():
        cur = 1
        for _ in range(count):
            cur = ((cur * p) + 1) % MOD
        res = res * cur % MOD
    return res


def sumOfFactorsOfFactors(counter: "Counter[int]") -> int:
    """返回所有`约数的约数之和`, counter 为这个数的所有质数因子分解."""
    res = 1
    for p, count in counter.items():
        inv = pow(p - 1, MOD - 2, MOD)
        tmp = inv * (inv * (pow(p, count + 2, MOD) - 1) - (count + 2)) % MOD
        res = res * tmp % MOD
    return res


def countOfFactors(counter: "Counter[int]") -> int:
    """返回约数个数, counter 为这个数的所有质数因子分解."""
    res = 1
    for count in counter.values():
        res *= count + 1
        res %= MOD
    return res


if __name__ == "__main__":
    n = int(input())
    counter = Counter()
    for _ in range(n):
        counter += getPrimeFactors(int(input()))
    print(sumOfFactors(counter))
