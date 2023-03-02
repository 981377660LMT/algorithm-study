# 累乘和 (幂级数求和/等比数列求和)


from math import log2
from typing import Tuple

MOD = int(1e9 + 7)


def powerSum(x: int, n: int) -> Tuple[int, int]:
    """O(logn) 求 (x^n, x^0 + ... + x^(n - 1))  n >= 1"""
    sum_, p = 1, x  # res = x^0 + ... + x^(len - 1), p = x^len
    start = int(log2(n)) - 1
    for d in range(start, -1, -1):
        sum_ *= p + 1
        p *= p
        if (n >> d) & 1:
            sum_ += p
            p *= x
        sum_ %= MOD
        p %= MOD
    return p, sum_


print(powerSum(3, 100000))  # (1024, 2047)
