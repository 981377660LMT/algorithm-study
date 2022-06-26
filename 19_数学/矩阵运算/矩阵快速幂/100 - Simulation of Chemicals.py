# 試験管に 3 種類の物質 A, B, C がそれぞれ a,b,c グラム入っているとき、次の 1 秒間で以下のような変化が起こります。
# 1秒后物品的量变为
# a(1-Xi) + bYi
# b(1-Yi) + cZi
# c(1-Zi) + aXi
# 开始时每个物品都是1g 求Ti秒后的物品量
import sys
from typing import List

Matrix = List[List[int]]


def matqpow(base: List[List[int]], exp: int) -> List[List[int]]:
    """矩阵快速幂 不取模"""

    def mul(m1: Matrix, m2: Matrix) -> Matrix:
        """矩阵相乘"""
        ROW, COL = len(m1), len(m2[0])
        res = [[0] * COL for _ in range(ROW)]
        for r, ROW in enumerate(m1):
            for c in range(COL):
                for k, v in enumerate(ROW):
                    res[r][c] += v * m2[k][c]
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


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)

q = int(input())
for _ in range(q):
    x, y, z, t = map(float, input().split())

    T = [[1 - x, y, 0], [0, 1 - y, z], [x, 0, 1 - z]]
    resT = matqpow(T, int(t))
    res1 = sum(resT[0])
    res2 = sum(resT[1])
    res3 = sum(resT[2])
    print(res1, res2, res3, sep=' ')

