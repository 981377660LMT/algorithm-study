# 无向图三角形顶点乘积之和
# n,m<=1e5
# Enumerate Triangles
from typing import List


MOD = 998244353


def enumerate_triangles(adjList: List[List[int]], values: List[int]) -> int:
    n = len(adjList)
    deg = [len(e) for e in adjList]
    G = [[] for _ in range(n)]
    for i in range(n):
        for j in adjList[i]:
            if deg[i] > deg[j]:
                G[i].append(j)
            elif deg[i] == deg[j] and i > j:
                G[i].append(j)

    res = 0
    visited = [False] * n
    for x in range(n):
        for z in G[x]:
            visited[z] = True
        for y in G[x]:
            for z in G[y]:
                if visited[z]:
                    res += (values[x] * values[y] % MOD) * values[z] % MOD
                    if res >= MOD:
                        res -= MOD
        for z in G[x]:
            visited[z] = False
    return res


n, m = map(int, input().split())
values = list(map(int, input().split()))
adjList = [[] for _ in range(n)]
for _ in range(m):
    u, v = map(int, input().split())
    adjList[u].append(v)
    adjList[v].append(u)

print(enumerate_triangles(adjList, values))
