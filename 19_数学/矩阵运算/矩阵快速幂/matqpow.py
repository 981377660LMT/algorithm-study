# !时间复杂度O(n^3logk)

from typing import List


Matrix = List[List[int]]


def mul(m1: Matrix, m2: Matrix, mod: int) -> Matrix:
    """矩阵相乘"""
    ROW, COL = len(m1), len(m2[0])

    res = [[0] * COL for _ in range(ROW)]
    for r in range(ROW):
        for c in range(COL):
            for i in range(ROW):
                res[r][c] = (res[r][c] + m1[r][i] * m2[i][c]) % mod

    return res


def matqpow1(base: Matrix, exp: int, mod: int) -> Matrix:
    """矩阵快速幂"""

    def inner(base: Matrix, exp: int, mod: int) -> Matrix:
        ROW, COL = len(base), len(base[0])
        res = [[0] * COL for _ in range(ROW)]
        for r in range(ROW):
            res[r][r] = 1

        bit = 0
        while exp:
            if exp & 1:
                res = mul(res, pow2[bit], mod)
            exp //= 2
            bit += 1
            if bit == len(pow2):
                pow2.append(mul(pow2[-1], pow2[-1], mod))
        return res

    pow2 = [base]
    return inner(base, exp, mod)


######################################################################
# numpy的矩阵快速幂

import numpy as np


def matqpow2(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
    """np矩阵快速幂"""

    def inner(base: np.ndarray, exp: int, mod: int) -> np.ndarray:
        res = np.eye(*base.shape, dtype=np.uint64)

        bit = 0
        while exp:
            if exp & 1:
                res = (res @ pow2[bit]) % mod
            exp //= 2
            bit += 1
            if bit == len(pow2):
                pow2.append((pow2[-1] @ pow2[-1]) % mod)
        return res

    pow2 = [base]
    return inner(base, exp, mod)


if __name__ == "__main__":
    n = 876543210987654321
    MOD = int(1e9 + 7)

    res = [[2], [1], [1]]  # 初始状态
    T = [[1, 1, 1], [1, 0, 0], [0, 1, 0]]
    resT = matqpow1(T, n - 3, MOD)
    res = mul(resT, res, MOD)
    assert res[0][0] == 639479200

    res = [[2], [1], [1]]  # 初始状态
    T = np.array([[1, 1, 1], [1, 0, 0], [0, 1, 0]], np.uint64)
    resT = matqpow2(T, n - 3, MOD)
    res = (resT @ res) % MOD
    assert res[0][0] == 639479200
