from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List

INF = int(1e20)

# 题目里V,E<=1e5，终点站个数<=10


def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int) -> List[int]:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        _, cur = heappop(pq)
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))

    return dist


n, m, k, s = map(int, input().split())
s -= 1

adjMap = defaultdict(lambda: defaultdict(lambda: INF))
for _ in range(m):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u][v] = min(adjMap[u][v], w)

for _ in range(k):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u][v] = min(adjMap[u][v], w)
    adjMap[v][u] = min(adjMap[v][u], w)

cost1, cost2, _ = map(int, input().split())
tos = list(map(int, input().split()))
tos = [n - 1 for n in tos]

distMap = defaultdict(list)
for target in set(tos) | {s}:
    distMap[target] = dijkstra(n, adjMap, target)

res = 0
roads = [s] + tos + [s]
for i, (cur, next) in enumerate(zip(roads, roads[1:])):
    weight = distMap[cur][next]
    res += distMap[cur][next]
    if i != len(roads) - 2:
        res += cost1 if res & 1 else cost2


print(res)
