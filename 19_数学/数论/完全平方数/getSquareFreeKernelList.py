from typing import List
from math import isqrt


def getSquareFreeKernelList(maxN: int) -> List[int]:
    maxN += 1
    res = [0] * maxN
    for i in range(1, maxN):
        if res[i] == 0:
            for j in range(1, isqrt(maxN // i) + 1):
                res[i * j * j] = i
    return res
