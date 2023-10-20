# 求前n项和

from math import log2
from typing import Tuple, List

MOD = int(1e9 + 7)


def powerSum(x: int, n: int, mod: int) -> Tuple[int, int]:
    """
    O(logn) 求 (x^n, x^0 + ... + x^(n - 1)) 模 mod

    x,n,mod >= 1
    """
    if mod == 1:
        return 0, 0
    sum_, p = 1, x  # res = x^0 + ... + x^(n - 1), p = x^n
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


# 等比数列求和
def powerSum2(n: int, a0: int, q: int) -> List[int]:
    """等比数列前n项和"""
    res = [a0]
    curSum, curItem = a0, a0
    for _ in range(n - 1):
        curItem *= q
        curSum += curItem
        res.append(curSum)
    return res


print(powerSum(2, 50000000, MOD))
print(powerSum2(50000000, 1, 2))
