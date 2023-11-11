"""
图的邻接矩阵幂的意义，数学期望，矩阵快速幂

给定一个n个点m条边的无向图,每个点有一个权值
!执行k次操作,每次随机选择一条边,将两个端点的权值换为两个端点的权值的平均值
求操作k次后每个点的权值的期望 (mod 1e9+7)
n<=100 k<=1e9
时间复杂度O(n^3logk)

Solution:
https://www.cnblogs.com/martian148/p/15531945.html
分别考虑 ai 对自己的贡献和其他对 ai 的贡献
trans[i][j] 表示 j 对 i 的贡献
在trans矩阵中,`对角线上的元素`表示自己对自己的贡献,`非对角线上的元素`表示和自己连边的点对自己的贡献
!对于顶点v,如果没有选中与v相连的边,那么v的权值不变  => 自己对自己贡献为 1 - deg[v]/(总边数*2)
!如果选中了与v相连的边,那么v的权值变为与v相连的点的权值的平均值 => 别人对自己贡献为 adjMatrix[i][j]/(总边数*2)
"""

import sys
from typing import List, Tuple
from matqpow import matqpow1, matmul

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def graphSmoothing(n: int, edges: List[Tuple[int, int]], values: List[int], k: int) -> List[int]:
    adjMatrix = [[0] * n for _ in range(n)]  # 邻接矩阵,1表示有边，0表示无边
    deg = [0] * n
    for u, v in edges:
        adjMatrix[u][v] += 1
        adjMatrix[v][u] += 1
        deg[u] += 1
        deg[v] += 1

    inv = pow(2 * len(edges), MOD - 2, MOD)
    trans = [[0] * n for _ in range(n)]
    for i in range(n):
        for j in range(n):
            if i != j:
                trans[i][j] = adjMatrix[i][j] * inv % MOD  # 别人对自己的影响
            else:  # 自己对自己的影响
                trans[i][j] = (1 - deg[i] * inv) % MOD

    trans = matqpow1(trans, k, MOD)
    res = [[num] for num in values]
    res = matmul(trans, res, MOD)
    return [row[0] for row in res]


if __name__ == "__main__":
    n, m, k = map(int, input().split())
    values = list(map(int, input().split()))
    edges = []
    for _ in range(m):
        u, v = map(int, input().split())
        edges.append((u - 1, v - 1))

    res = graphSmoothing(n, edges, values, k)
    print(*res, sep="\n")
