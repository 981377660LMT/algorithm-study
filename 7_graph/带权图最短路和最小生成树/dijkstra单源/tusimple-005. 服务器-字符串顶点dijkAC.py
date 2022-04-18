from collections import defaultdict
from functools import lru_cache
from heapq import heappop, heappush
from typing import DefaultDict


def dijkstra(adjMap: DefaultDict[str, DefaultDict[str, int]], start: str, end: str) -> int:
    """时间复杂度O((V+E)logV)"""

    @lru_cache(None)
    def inner(start: str, end: str) -> int:
        pq = [(0, start)]
        dist = defaultdict(lambda: int(1e20))
        dist[start] = 0
        while pq:
            curDist, cur = heappop(pq)
            if cur == end:
                return curDist
            for next in adjMap[cur]:
                if dist[next] > dist[cur] + adjMap[cur][next]:
                    dist[next] = dist[cur] + adjMap[cur][next]
                    heappush(pq, (dist[next], next))

        return int(1e20)

    return inner(start, end)


# 2 <= n <= 1000,1 <= m <= 1000,1 <= q <= 1000
n, m = map(int, input().split())

adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))

for _ in range(m):
    u, v, w = input().split()
    w = int(w)
    adjMap[u][v] = min(adjMap[u][v], w)


q = int(input())
for _ in range(q):
    u, v = input().split()
    res = dijkstra(adjMap, u, v)
    print(res if res < int(1e19) else 'INF')
