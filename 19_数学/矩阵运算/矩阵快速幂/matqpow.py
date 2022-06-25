from typing import List
import numpy as np


Matrix = List[List[int]]


# 普通的矩阵快速幂 (注意如果要js写的话 需要bigint)
def matqpow1(base: List[List[int]], exp: int, mod=int(1e9 + 7)) -> List[List[int]]:
    """矩阵快速幂"""

    def mul(m1: Matrix, m2: Matrix) -> Matrix:
        """矩阵相乘"""
        ROW, COL = len(m1), len(m2[0])
        res = [[0] * COL for _ in range(ROW)]
        for r, ROW in enumerate(m1):
            for c in range(COL):
                for k, v in enumerate(ROW):
                    res[r][c] += v * m2[k][c]
                    res[r][c] %= mod
        return res

    ROW, COL = len(base), len(base[0])
    res = [[0] * COL for _ in range(ROW)]
    for r in range(ROW):
        res[r][r] = 1

    while exp:
        if exp & 1:
            res = mul(res, base)
        exp >>= 1
        base = mul(base, base)
    return res


# numpy的矩阵快速幂 (注意如果要js写的话 需要bigint)
def matqpow2(base: Matrix, exp: int, mod=int(1e9 + 7)) -> Matrix:
    """矩阵快速幂np版"""

    def mul(m1: np.ndarray, m2: np.ndarray) -> np.ndarray:
        """矩阵相乘"""
        return np.dot(m1, m2) % mod

    ROW = len(base)
    res = np.eye(ROW, dtype=np.int64)
    npBase = np.array(base, dtype=np.int64)

    while exp:
        if exp & 1:
            res = mul(res, npBase)
        exp >>= 1
        npBase = mul(npBase, npBase)
    return res.tolist()


if __name__ == '__main__':
    n = 876543210987654321
    MOD = int(1e9 + 7)
    T = [[1, 1, 1], [1, 0, 0], [0, 1, 0]]
    resT = matqpow2(T, n - 3, MOD)
    a3, a2, a1 = 2, 1, 1
    res = (resT[0][0] * a3 + resT[0][1] * a2 + resT[0][2] * a1) % MOD
    assert res == 639479200
