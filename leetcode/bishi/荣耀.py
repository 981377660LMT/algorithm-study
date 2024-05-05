# n<=1000
# q<=1000

from collections import defaultdict
from heapq import heappop, heappush
import sys
from typing import Mapping

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


n, q = map(int, input().split())
adjMap = defaultdict(lambda: defaultdict(lambda: INF))
edges = []
for _ in range(2 * n - 2):
    u, v, target = map(int, input().split())
    u, v = u - 1, v - 1
    adjMap[u][v] = min(adjMap[u][v], target)
    edges.append((u, v))


def dijkstra(n: int, adjMap: Mapping[int, Mapping[int, int]], start: int, end: int) -> int:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if cur == end:
            return curDist
        if curDist > dist[cur]:
            continue
        for next, nextDist in adjMap[cur].items():
            if curDist + nextDist < dist[next]:
                dist[next] = curDist + nextDist
                heappush(pq, (dist[next], next))
    return INF


# !实时查询两个景点间的最短距离
for _ in range(q):
    type, *rest = map(int, input().split())
    if type == 1:
        ei, target = rest  # !将第i条车道的长度调整为w
        ei -= 1
        u, v = edges[ei]
        adjMap[u][v] = target
    else:
        u, v = rest  # !打印景点u到景点v的最短路径
        u, v = u - 1, v - 1
        print(dijkstra(n, adjMap, u, v))
