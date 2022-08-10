# h<=1e9
# r<=16
# 房子有h层 每层共有r个房间
# 如果gij=1 那么房间i和房间j存在一条无向边连通
# 层之间可以坐梯子下一层楼(下楼后所在列不变)
# !求h层的1号房间到1层的1号房间的路径数 不能重复进入一个房间

# !行间矩阵快速幂(col^3*logrow) + 行内状压dp预处理(col^3*2^col)
# 1. 行内状压dp求i->j的方案数
# 2. 行间矩阵快速幂 注意到`dp转移类似于矩阵乘法` 所以用矩阵快速幂优化


import sys
from typing import List


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)
INF = int(4e18)


Matrix = List[List[int]]


def mul(m1: Matrix, m2: Matrix, mod: int) -> Matrix:
    ROW, COL = len(m1), len(m2[0])
    res = [[0] * COL for _ in range(ROW)]
    for r in range(ROW):
        for c in range(COL):
            for i in range(ROW):
                res[r][c] += m1[r][i] * m2[i][c]
                res[r][c] %= mod
    return res


def matqpow1(base: Matrix, exp: int, mod: int) -> Matrix:
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


#########################################################################
ROW, COL = map(int, input().split())
adjMatrix = []
for _ in range(COL):
    row = list(map(int, input().split()))
    adjMatrix.append(row)


memo = [-1] * (COL * COL * 1 << COL)


def dfs(cur: int, target: int, visited: int) -> int:
    if cur == target:
        return 1
    hash_ = visited * COL * COL + cur * COL + target
    if memo[hash_] != -1:
        return memo[hash_]
    res = 0
    for next in range(COL):
        if (visited & (1 << next)) == 0 and adjMatrix[cur][next] == 1:
            res += dfs(next, target, visited | (1 << next))
    memo[hash_] = res
    return res


ways = [[1] * COL for _ in range(COL)]  # 同一层 i -> j 的方案数
for i in range(COL):
    for j in range(i + 1, COL):
        ways[i][j] = ways[j][i] = dfs(i, j, 1 << i)


# 两层时: ∑ (ways[1][j]*ways[j][1])
# 三层时: ∑ (ways[1][j]*ways[j][k]*ways[k][1])
# !和矩阵相乘的含义一样


T = matqpow1(ways, ROW, MOD)  # 转移矩阵
init = [[0] for _ in range(COL)]
init[0][0] = 1
res = mul(T, init, MOD)
print(res[0][0])
