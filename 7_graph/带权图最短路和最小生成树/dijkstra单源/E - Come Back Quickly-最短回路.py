# E - Come Back Quickly
# 有重边有自环的有向图(无负环) 求从每个点出发至少经过1条边回到起点的最短回路
# 如果不存在这样的路径输出-1
# n<=2000 m<=2000

# !先将所有与start点相连的点加入队列再进行dijkstra

from typing import List, Sequence, Tuple
from heapq import heappop, heappush

INF = int(1e18)


def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int) -> List[int]:
    dist = [INF] * n
    pq = []
    for next, weight in adjList[start]:  # 将所有与start点相连的点加入队列
        dist[next] = min(dist[next], weight)
        heappush(pq, (dist[next], next))

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
    return dist


def comeBackQuickly(n: int, edges: List[Tuple[int, int, int]]) -> List[int]:
    adjList = [[] for _ in range(n)]
    for u, v, w in edges:
        adjList[u].append((v, w))
    return [dijkstra(n, adjList, i)[i] for i in range(n)]


n, m = map(int, input().split())
edges = []
for _ in range(m):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    edges.append((u, v, w))
res = comeBackQuickly(n, edges)
for dist in res:
    print(dist if dist < INF else -1)
