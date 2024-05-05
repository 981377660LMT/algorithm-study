# z0_可能在最短路径上的边
# https://leetcode.cn/problems/find-edges-in-shortest-paths/description/
# https://zhuanlan.zhihu.com/p/166173790


from typing import List, Sequence, Tuple
from heapq import heappop, heappush

INF = int(1e18)


def edgesMayOnShortestPath(
    n: int, edges: List[Tuple[int, int, int]], start: int, target: int
) -> List[bool]:
    """
    给定一个n个点m条边的无向带权图.
    对于每条边(u, v, w), 判断是否可能在从start到target的最短路径上.
    正反跑一次dijkstra, 然后判断是否满足最短路径的条件：
    dist1[u] + w + dist2[v] == dist1[target] or dist1[v] + w + dist2[u] == dist1[target]
    """
    adjList = [[] for _ in range(n)]
    for _, (u, v, w) in enumerate(edges):
        adjList[u].append((v, w))
        adjList[v].append((u, w))
    dist1 = dijkstra(n, adjList, start)
    dist2 = dijkstra(n, adjList, target)
    res = [False] * len(edges)
    for i, (u, v, w) in enumerate(edges):
        if dist1[u] + w + dist2[v] == dist1[target] or dist1[v] + w + dist2[u] == dist1[target]:
            res[i] = True
    return res


def dijkstra(n: int, adjList: Sequence[Sequence[Tuple[int, int]]], start: int) -> List[int]:
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


class Solution:
    def findAnswer(self, n: int, edges: List[Tuple[int, int, int]]) -> List[bool]:
        return edgesMayOnShortestPath(n, edges, 0, n - 1)
