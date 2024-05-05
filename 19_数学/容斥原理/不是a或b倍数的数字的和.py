"""求倍数个数"""
# D - FizzBuzz Sum Hard-不是a或b倍数的数字的和
# 给定n, a,b找出[1,n]内不是a或b倍数的数字的和。
#
# !简单的容斥原理，扣掉a, b 的倍数，加上lcm(a,b)的倍数。


from math import gcd
from typing import Tuple


def calMul1(base: int, upper: int) -> Tuple[int, int]:
    """返回 [1,upper] 中 base 的 (倍数的个数,倍数的和)"""
    count = upper // base
    last = base + (count - 1) * base
    sum_ = (base + last) * count // 2
    return count, sum_


def calMul2(base: int, lower: int, upper: int) -> Tuple[int, int]:
    """返回 [lower,upper] 中 base 的 (倍数的个数,倍数的和)"""
    count1, sum1 = calMul1(base, upper)
    count2, sum2 = calMul1(base, lower - 1)
    return count1 - count2, sum1 - sum2


import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    n, a, b = map(int, input().split())
    allSum = (1 + n) * n // 2
    print(allSum - calMul2(a, 1, n)[1] - calMul2(b, 1, n)[1] + calMul2(a * b // gcd(a, b), 1, n)[1])


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
