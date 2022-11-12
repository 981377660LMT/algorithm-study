# 有两种形状的瓷砖：一种是 2x1 的多米诺形，
# 另一种是形如 "L" 的托米诺形。两种形状都可以旋转。
# !给定 N 的值，有多少种方法可以平铺 2 x N 的面板？返回值 mod 10^9 + 7。
# 多米诺和托米诺平铺(结论题)
# !一维DP状态转移方程：dp[i] = 2 * dp[i - 1] + dp[i - 3]

# !即:
# [ai  ]   =    [2,0,1,0]  * [ai-1]
# [ai-1]        [1,0,0,0]    [ai-2]
# [ai-2]        [0,1,0,0]    [ai-3]
# [ai-3]        [0,0,1,0]    [ai-4]
# n<=3时 取init
# n>3 时，转移n-3次

MOD = int(1e9 + 7)


class Solution:
    def numTilings(self, n: int) -> int:
        init = [[5], [2], [1], [0]]
        if n <= 3:
            return init[~n][0]
        T = [[2, 0, 1, 0], [1, 0, 0, 0], [0, 1, 0, 0], [0, 0, 1, 0]]
        resT = matqpow1(T, n - 3, MOD)
        res = mul(resT, init, MOD)
        return res[0][0]


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
