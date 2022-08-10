# 求有向图中长为k的路径数
# n<=50 k<=1e18

# dp[i][j][l] 表示从i到j的长为l的路径条数
# dp[i][j][l] = ∑(dp[i][k][l-1]*adj[k][j])
# !注意到`dp转移类似于矩阵乘法` 所以用矩阵快速幂优化
# !时间复杂度O(n^3logk)

import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)
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
        exp //= 2
        base = mul(base, base, mod)
    return res


#############################################################

n, k = map(int, input().split())
adjMatrix = []
for _ in range(n):
    adjMatrix.append(list(map(int, input().split())))

T = matqpow1(adjMatrix, k, MOD)  # dp 矩阵
res = 0
for r in range(n):
    for c in range(n):
        res += T[r][c]
        res %= MOD
print(res)
