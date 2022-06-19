# !所有因数
from math import isqrt
from typing import List


def getFactors(n: int) -> List[int]:
    """返回 n 的所有因数"""
    upper = isqrt(n) + 1
    small, big = [], []

    for i in range(1, upper):
        if n % i == 0:
            small.append(i)
            big.append(n // i)

    if small[-1] == big[-1]:
        small.pop()
    return small + big[::-1]


class Solution:
    def isThree(self, n: int) -> bool:
        return len(getFactors(n)) == 3

