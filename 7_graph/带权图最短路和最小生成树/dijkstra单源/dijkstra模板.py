"""dijk模板"""

from heapq import heappop, heappush
from typing import List, Sequence, Tuple

INF = int(1e18)


def dijkstra1(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int) -> List[int]:
    """dijkstra求出起点到各点的最短距离 时间复杂度O((V+E)logV)"""
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]

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


# dijkstra求出一条路径
def dijkstra2(
    n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int, end: int
) -> Tuple[int, List[int]]:
    """dijkstra求出起点到end的(最短距离,路径) 时间复杂度O((V+E)logV)"""
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    pre = [-1] * n  # 记录每个点的前驱

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (dist[next], next))
                pre[next] = cur

    if dist[end] == INF:
        return INF, []

    path = []
    cur = end
    while pre[cur] != -1:
        path.append(cur)
        cur = pre[cur]
    path.append(start)
    return dist[end], path[::-1]


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")

if __name__ == "__main__":
    n, m, start, end = map(int, input().split())
    adjList = [[] for _ in range(n)]
    for _ in range(m):
        u, v, w = map(int, input().split())
        adjList[u].append((v, w))

    dist, path = dijkstra2(n, adjList, start, end)
    if dist == INF:
        print(-1)
        exit(0)
    print(dist, len(path) - 1)
    for a, b in zip(path, path[1:]):
        print(a, b)
