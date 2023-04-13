# https://maspypy.github.io/library/seq/find_linear_rec.hpp
from typing import List

MOD = int(1e9 + 7)


def findLinearRec(nums: List[int]) -> int:
    """寻找数列的线性递推式系数."""
    n = len(nums)
    B, C = [1], [1]
    l, m = 0, 1
    p = 1
    for i in range(n):
        d = nums[i]
        for j in range(1, l + 1):
            d += C[j] * nums[i - j]
        d %= MOD
        if d == 0:
            m += 1
            continue
        tmp = C[:]
        q = d // p
        if len(C) < len(B) + m:
            C += [0] * (len(B) + m - len(C))
        for j in range(len(B)):
            C[j + m] = (C[j + m] - q * B[j]) % MOD
        if l + l <= i:
            B = tmp
            l = i + 1 - l
            m = 1
            p = d
        else:
            m += 1
    return C


print(findLinearRec([1, 1, 2, 3, 5]))
