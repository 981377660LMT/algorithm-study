#  https://yukicoder.me/problems/no/1320
#  n,m<=2000,wi<=1e9
# !枚举每条边删除：`O(E*(V+E)*logV)` 找最小环,不存在返回INF

from collections import deque
from heapq import heappop, heappush
from typing import List, Tuple

INF = int(1e18)


def minCostCycle(n: int, edges: List[Tuple[int, int, int]], directed: bool) -> int:
    """O(E*(V+E)*logV)求最小环权值和,不存在返回INF."""
    adjList = [[] for _ in range(n)]
    maxWeight = 0
    for u, v, w in edges:
        if w > maxWeight:
            maxWeight = w
        adjList[u].append((v, w))
        if not directed:
            adjList[v].append((u, w))
    res = INF
    for i in range(len(edges)):  # remove edge i
        from_, to, weight = edges[i]  # 注意有向图时最短路需要to->from
        dist = _bfs01(adjList, to, from_) if maxWeight <= 1 else _dijkstra(adjList, to, from_)
        cand = weight + dist
        if cand < res:
            res = cand
    return res


def _bfs01(adjList: List[List[Tuple[int, int]]], start: int, target: int) -> int:
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    while queue:
        cur = queue.popleft()
        if cur == target:
            return dist[cur]
        for next, weight in adjList[cur]:
            if (cur == start and next == target) or (cur == target and next == start):
                continue
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                if weight == 0:
                    queue.appendleft(next)
                else:
                    queue.append(next)
    return INF


def _dijkstra(adjList: List[List[Tuple[int, int]]], start: int, target: int) -> int:
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if cur == target:
            return curDist
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            if (cur == start and next == target) or (cur == target and next == start):
                continue
            cand = curDist + weight
            if cand < dist[next]:
                dist[next] = cand
                heappush(pq, (cand, next))
    return INF


if __name__ == "__main__":
    directed = int(input())
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        u -= 1
        v -= 1
        edges.append((u, v, w))
    res = minCostCycle(n, edges, directed == 1)
    if res == INF:
        res = -1
    print(res)

    class Solution:
        def findShortestCycle(self, n: int, edges: List[List[int]]) -> int:
            newEdges = [(u, v, 1) for u, v in edges]
            res = minCostCycle(n, newEdges, False)
            return res if res != INF else -1
