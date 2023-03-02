# 求前n项和
from typing import List


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


print(powerSum(2, 50000000))
print(powerSum2(50000000, 1, 2))
