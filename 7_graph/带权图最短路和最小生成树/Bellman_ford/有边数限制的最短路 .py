# 给定一个 n 个点 m 条边的有向图，图中可能存在重边和自环， 边权可能为负数。
# 请你求出从 1 号点到 n 号点的最多经过 k 条边的最短距离，如果无法从 1 号点走到 n 号点，输出 impossible。
#  边权可能为负数。
V, E, limit = map(int, input().split())
dist = [int(1e20)] * (V + 1)
dist[1] = 0

edges = [tuple(map(int, input().split())) for _ in range(E)]

for _ in range(limit):  # !中转点数i=0,1,...,limit-1
    preDist = dist[:]
    for u, v, w in edges:
        if preDist[u] + w < dist[v]:
            dist[v] = preDist[u] + w

print("impossible" if dist[V] > int(1e19) else (dist[V]))
