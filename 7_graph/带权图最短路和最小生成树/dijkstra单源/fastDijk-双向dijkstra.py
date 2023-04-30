#
#  Dijkstra's Single Source Shortest Path
#
#  Description:
#
#    Dijkstra algorithm finds a single source shortest path on
#    a nonnegative weighted graph.
#    It implements two algorithms with two data structures.
#    (1) standard dijkstra with standard heap
#    (2) bidirectional dijkstra with standard heap
#    (3) standard dijkstra with radix heap
#    (4) bidirectional dijkstra with radix heap
#
#    For a simple test (see the code), we observe that
#      Binomial,Unidirectional  14.74[s]
#      Binomial Bidirectional    0.52[s]
#      Radix,Unidirectional      4.90[s]
#      Radix,Bidirectional       0.31[s]
#
#  Complexity:
#
#    O(m log n) for binomial heap
#    O(m) for radix heap
#
#  Verify:
#
#    SPOJ_SHPATH  (bidirectional search)
#


from heapq import heappop, heappush
from typing import List, Sequence, Tuple


def fastDijkstra(
    n: int,
    adjList: Sequence[Sequence[Tuple[int, int]]],
    rAdjList: Sequence[Sequence[Tuple[int, int]]],
    start: int,
    end: int,
) -> int:
    """双向dijkstra求start到end的最短路."""
    if start == end:
        return 0
    dist = [-1] * n
    dist[start] = 0
    drev = [-1] * n
    drev[end] = 0
    pq1, pq2 = [(0, start)], [(0, end)]
    res = -1
    while pq1 and pq2:
        d1, d2 = pq1[0][0], pq2[0][0]
        if res >= 0 and d1 + d2 >= res:
            break
        if d1 <= d2:
            d, u = heappop(pq1)
            if dist[u] > d:
                continue
            for v, w in adjList[u]:
                cand = dist[u] + w
                if dist[v] >= 0 and dist[v] <= cand:
                    continue
                dist[v] = cand
                heappush(pq1, (dist[v], v))
                if drev[v] >= 0:
                    nu = dist[v] + drev[v]
                    if res < 0 or res > nu:
                        res = nu
        else:
            d, u = heappop(pq2)
            if drev[u] > d:
                continue
            for v, w in rAdjList[u]:
                cand = drev[u] + w
                if drev[v] >= 0 and drev[v] <= cand:
                    continue
                drev[v] = cand
                heappush(pq2, (drev[v], v))
                if dist[v] >= 0:
                    nu = dist[v] + drev[v]
                    if res < 0 or res > nu:
                        res = nu
    return res


if __name__ == "__main__":

    class Graph:
        def __init__(self, n: int, edges: List[List[int]]):
            adjList = [[] for _ in range(n)]
            rAdjList = [[] for _ in range(n)]
            for u, v, w in edges:
                adjList[u].append((v, w))
                rAdjList[v].append((u, w))
            self.adjList = adjList
            self.rAdjList = rAdjList

        def addEdge(self, edge: List[int]) -> None:
            u, v, w = edge
            self.adjList[u].append((v, w))
            self.rAdjList[v].append((u, w))

        def shortestPath(self, node1: int, node2: int) -> int:
            return fastDijkstra(len(self.adjList), self.adjList, self.rAdjList, node1, node2)
