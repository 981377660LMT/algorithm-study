"""
边长为d的正n边形
每个点染成黑色或白色
每条边的白色个数相等
求染色方法数

n<=1e12
d<=1e4
时间复杂度dlogn
"""

import sys
import os

from functools import lru_cache
from typing import List


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = 998244353


Matrix = List[List[int]]


# 普通的矩阵快速幂 (注意如果要js写的话 需要bigint)
def matqpow1(base: List[List[int]], exp: int, mod=int(1e9 + 7)) -> List[List[int]]:
    """矩阵快速幂"""

    def mul(m1: Matrix, m2: Matrix) -> Matrix:
        """矩阵相乘"""
        ROW, COL = len(m1), len(m2[0])
        res = [[0] * COL for _ in range(ROW)]
        for r, row in enumerate(m1):
            for c in range(COL):
                for k, v in enumerate(row):
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
        exp //= 2
        base = mul(base, base)
    return res


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


@lru_cache(None)
def C(n: int, k: int) -> int:
    if n < 0 or k < 0 or n < k:
        return 0
    return ((fac(n) * ifac(k)) % MOD * ifac(n - k)) % MOD


def main() -> None:
    n, d = map(int, input().split())
    # 枚举最后每条边上多少个白色点
    # !dp[i][s1,s2] 表示前i条边 最后的端点放黑色还是白色 注意每条边有d+1个点
    # !dp[i][0]=C(d-1,k)*dp[i-1][0]+C(d-1,k-1)*dp[i-1][1]
    # !dp[i][1]=C(d-1,k-1)*dp[i-1][0]+C(d-1,k-2)*dp[i-1][1]
    # !dp转移式不依赖于i 求第n项可以用矩阵快速幂
    res = 0
    for k in range(d + 2):
        T = [[C(d - 1, k), C(d - 1, k - 1)], [C(d - 1, k - 1), C(d - 1, k - 2)]]
        resT = matqpow1(T, n, MOD)
        res += resT[0][0] + resT[1][1]  # 最开始放黑 最后黑结尾 + 最开始放白 最后白结尾
        res %= MOD
    print(res)


if os.environ.get("USERNAME", "") == "caomeinaixi":
    while True:
        try:
            main()
        except EOFError:
            break
else:
    main()
