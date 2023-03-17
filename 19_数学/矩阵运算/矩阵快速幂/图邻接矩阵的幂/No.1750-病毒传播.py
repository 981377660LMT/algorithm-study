# No.1750-病毒传播
# 给定一个图
# 初始时0点被感染
# !第i天每个都市的感染人数为前一天`邻接的所有城市的感染人数之和`
# 求第k天0号城市的感染人数 模 998244353

# n<=100 k<=1e9


import numpy as np


def matqpow2(base: "np.ndarray", exp: int, mod: int) -> "np.ndarray":
    """np矩阵快速幂"""
    res = np.eye(*base.shape, dtype=np.int128)
    while exp:
        if exp & 1:
            res = (res @ base) % mod
        base = (base @ base) % mod
        exp >>= 1
    return res


n, m, k = map(int, input().split())
adjMatrix = [[0] * n for _ in range(n)]

for _ in range(m):
    u, v = map(int, input().split())
    adjMatrix[u][v] = 1
    adjMatrix[v][u] = 1

T = matqpow2(np.array(adjMatrix, np.int128), k, 998244353)
print(T[0][0])
