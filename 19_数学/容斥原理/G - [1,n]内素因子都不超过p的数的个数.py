import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# 求[1,n]内素因子都不超过p的数的个数
# 容斥原理
from collections import Counter
from typing import List


def isPrime(n: int) -> bool:
    """判断n是否为质数"""
    if n < 2:
        return False
    for i in range(2, int(n**0.5) + 1):
        if n % i == 0:
            return False
    return True


if __name__ == "__main__":
    N, P = map(int, input().split())
    primes = [i for i in range(2, P + 1) if isPrime(i)]

    # 枚举素因子的个数???
    res = 0

    for pCount in range(60):
        ...
