# 寻找一般图最大独立集
# n<=40
# Maximum Independent Set
# 独立集：在一个图中，找到一个集合包含的所有点相互之间都不存在连边
# 最大独立集：在所有独立集中包含元素个数最多的独立集


from typing import List

# https://judge.yosupo.jp/submission/35335
def maximum_clique_BK(adjMatrix: List[List[int]]) -> List[int]:
    def dfs(R, P, cnt):
        nonlocal res, val
        if popcount(P) + cnt <= val:
            return  # 枝刈り
        if P == 0:
            if val < cnt:
                val = cnt
                res = R
            return
        PP = P & ~adjList[(P & -P).bit_length() - 1]
        while PP:
            v = PP & -PP
            vv = v.bit_length() - 1
            dfs(R | v, P & adjList[vv], cnt + 1)
            P ^= v
            PP ^= v

    n = len(adjMatrix)
    # assert n <= 63
    N = 1 << n
    adjList = [0] * n
    for a in range(n):
        for b in range(n):
            if adjMatrix[a][b] == 0 and a != b:
                adjList[a] ^= 1 << b
    res = val = 0
    popcount = popcountll if n <= 64 else popcount128
    dfs(R=0, P=N - 1, cnt=0)
    clique = [i for i, c in enumerate(bin(res)[2:][::-1]) if c == "1"]
    return clique


def popcountll(i):
    i -= (i >> 1) & 0x5555555555555555
    i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
    i = (i & 0x0F0F0F0F0F0F0F0F) + ((i >> 4) & 0x0F0F0F0F0F0F0F0F)
    i = (i >> 32) + (i & 0xFFFFFFFF)
    return ((i * 0x1010101) & 0xFFFFFFFF) >> 24


def popcount128(i):
    return popcountll(i >> 62) + popcountll(i & 0x3FFFFFFFFFFFFFFF)


n, m = map(int, input().split())
adjMatrix = [[0] * n for _ in range(n)]
for _ in range(m):
    u, v = map(int, input().split())
    adjMatrix[u][v] = adjMatrix[v][u] = 1
maxClique = maximum_clique_BK(adjMatrix)
print(len(maxClique))
print(*maxClique)
