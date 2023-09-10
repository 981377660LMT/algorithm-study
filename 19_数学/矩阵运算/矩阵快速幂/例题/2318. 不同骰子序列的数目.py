# 36*36维度的转移矩阵
# 时间复杂度O(logn*d^3)
from functools import lru_cache
from itertools import permutations
from math import gcd
import numpy as np

gcd = lru_cache(gcd)
MOD = int(1e9 + 7)


class Solution:
    def distinctSequences(self, n: int) -> int:
        if n <= 2:
            return [0, 6, 22][n]

        init = np.ones((36,), np.uint64)  # 初始状态
        T = np.zeros((36, 36), np.uint64)  # 转移矩阵
        for pre1, pre2, pre3 in permutations(range(1, 7), 3):
            if gcd(pre1, pre2) == 1 and gcd(pre2, pre3) == 1:
                T[(pre2 - 1) * 6 + (pre3 - 1)][(pre1 - 1) * 6 + (pre2 - 1)] = 1

        resT = matqpow2(T, n - 2, MOD)
        res = (resT @ init) % MOD
        return int(sum(res)) % MOD


def matqpow2(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
    """np矩阵快速幂"""

    base = base.copy()
    res = np.eye(*base.shape, dtype=np.uint64)

    while exp:
        if exp & 1:
            res = (res @ base) % mod
        exp //= 2
        base = (base @ base) % mod
    return res
