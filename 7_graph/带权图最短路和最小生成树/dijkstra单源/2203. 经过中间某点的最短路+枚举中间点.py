# 给一个无向带权图，求出从start到end，期间必须经过mid的最短距离
# 如果不存在，输出-1
# !min(x, a) + min(x, b) + min(x, c)

from typing import List, Mapping
from collections import defaultdict
from heapq import heappop, heappush


INF = int(1e18)


def dijkstra(n: int, adjMap: Mapping[int, Mapping[int, int]], start: int) -> List[int]:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next in adjMap[cur]:
            cand = dist[cur] + adjMap[cur][next]
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
    return dist


def solve(n: int, edges: List[List[int]], start: int, end: int, mid: int) -> int:
    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    for u, v, w in edges:
        adjMap[u][v] = min(adjMap[u][v], w)
        adjMap[v][u] = min(adjMap[v][u], w)

    dist1 = dijkstra(n, adjMap, start)
    dist2 = dijkstra(n, adjMap, mid)
    dist3 = dijkstra(n, adjMap, end)

    res = INF
    for i in range(n):
        cand = dist1[i] + dist2[i] + dist3[i]
        res = cand if cand < res else res

    return res if res < INF else -1
