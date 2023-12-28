# 累乘和 (幂级数求和/等比数列求和)


from math import log2
from typing import Tuple


def powerSum(x: int, n: int, mod: int) -> Tuple[int, int]:
    """
    O(logn) 求 (x^n, x^0 + ... + x^(n - 1)) 模 mod

    x,n,mod >= 1
    """
    if mod == 1:
        return 0, 0
    sum_, p = 1, x  # res = x^0 + ... + x^(len - 1), p = x^len
    start = int(log2(n)) - 1

    for d in range(start, -1, -1):
        sum_ *= p + 1
        p *= p
        if (n >> d) & 1:
            sum_ += p
            p *= x
        sum_ %= mod
        p %= mod
    return p, sum_


if __name__ == "__main__":
    # # https://atcoder.jp/contests/abc293/editorial/5955
    # print(powerSum(3, 100000, MOD))  # (1024, 2047)
    A, X, M = map(int, input().split())
    print(powerSum(A, X, M)[1])
