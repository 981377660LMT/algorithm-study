# 注意如果要js写的话 需要bigint

from typing import List


Matrix = List[List[int]]


def mul(m1: Matrix, m2: Matrix, mod: int) -> Matrix:
    """矩阵相乘"""
    ROW, COL = len(m1), len(m2[0])

    res = [[0] * COL for _ in range(ROW)]
    for r in range(ROW):
        for c in range(COL):
            for i in range(ROW):
                res[r][c] += m1[r][i] * m2[i][c]
                res[r][c] %= mod

    return res


# 普通的矩阵快速幂
def matqpow1(base: Matrix, exp: int, mod: int) -> Matrix:
    """矩阵快速幂"""

    ROW, COL = len(base), len(base[0])
    res = [[0] * COL for _ in range(ROW)]
    for r in range(ROW):
        res[r][r] = 1

    while exp:
        if exp & 1:
            res = mul(res, base, mod)
        exp >>= 1
        base = mul(base, base, mod)
    return res


######################################################################
# numpy的矩阵快速幂

import numpy as np


def matqpow2(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
    """np矩阵快速幂"""

    base = base.copy()
    res = np.eye(*base.shape, dtype=np.uint64)

    while exp:
        if exp & 1:
            res = (res @ base) % mod
        exp >>= 1
        base = (base @ base) % mod
    return res


if __name__ == "__main__":
    n = 876543210987654321
    MOD = int(1e9 + 7)

    res = [[2], [1], [1]]  # 初始状态
    trans = [[1, 1, 1], [1, 0, 0], [0, 1, 0]]
    resT = matqpow1(trans, n - 3, MOD)
    res = mul(resT, res, MOD)
    assert res[0][0] == 639479200

    res = [[2], [1], [1]]  # 初始状态
    trans = np.array([[1, 1, 1], [1, 0, 0], [0, 1, 0]], np.uint64)
    resT = matqpow2(trans, n - 3, MOD)
    res = (resT @ res) % MOD
    assert res[0][0] == 639479200
