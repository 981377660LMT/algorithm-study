"""
1<=n<=1e18
有多少个整数k=p*q^3 其中 p<q

枚举q
"""


from bisect import bisect_right
import sys
import os
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def getPrimes(n: int) -> List[int]:
    """筛法求小于等于n的素数"""
    isPrime = [True] * (n + 1)
    for num in range(2, n + 1):
        if isPrime[num]:
            for multi in range(num * num, n + 1, num):
                isPrime[multi] = False
    return [num for num in range(2, n + 1) if isPrime[num]]


primes = getPrimes(int(1e6 + 5))
ok = set(primes)


def main() -> None:
    # n = int(input())
    # res = 0
    # for q in primes:
    #     if q**3 > n:
    #         break
    #     upper = min(n // q**3, q - 1)
    #     # 小于等于upper和q-1的质数
    #     res += bisect_right(primes, upper)  # !也可以用前缀和求小于等于upper的质数多少个
    # print(res)

    n = int(input())
    res = 0
    preSum = [0, 0, 1]  # 小于等于i的质数个数
    for i in range(3, len(primes) + 1):
        preSum.append(preSum[-1] + int(i in ok))
    for q in primes:
        if q**3 > n:
            break
        upper = min(n // q**3, q - 1)
        res += preSum[upper]  # 前缀和求小于等于upper的质数多少个
    print(res)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
