# 求第N项模1e9+7的值
# a1=1 a2=1 a3=2 an=a(n-1)+a(n-2)+a(n-3)
# N<=1e18

# [an  ]     =  [1 1 1]   * [an-1]
# [an-1]        [1 0 0]     [an-2]
# [an-2]        [0 1 0]     [an-3]
import numpy as np


NPArray = np.ndarray


def matqpow2(base: NPArray, exp: int, mod: int) -> NPArray:
    """矩阵快速幂np版"""

    base = base.copy()
    res = np.eye(*base.shape, dtype=np.uint64)

    while exp:
        if exp & 1:
            res = (res @ base) % mod
        exp //= 2
        base = (base @ base) % mod
    return res


import sys

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

n = int(input())

if n <= 3:
    print(2 if n == 3 else 1)
    exit(0)

res = np.array([[2], [1], [1]], np.uint64)
T = np.array([[1, 1, 1], [1, 0, 0], [0, 1, 0]], np.uint64)
resT = matqpow2(T, n - 3, MOD)
res = (resT @ res) % MOD
print(int(res[0][0]))
