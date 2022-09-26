# bfs模板 类似于dijkstra

from collections import deque
from typing import List, Tuple

INF = int(1e18)


def bfs(start: int, adjList: List[List[int]]) -> List[int]:
    """时间复杂度O(V+E)"""
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    while queue:
        cur = queue.popleft()
        for next in adjList[cur]:
            cand = dist[cur] + 1
            if cand < dist[next]:
                dist[next] = cand
                queue.append(next)
    return dist


def bfs2(
    start: int, adjList: List[List[Tuple[int, int]]], banEdge: int
) -> Tuple[List[int], List[Tuple[int, int]]]:
    """bfs求最短路,并记录一条路径"""
    n = len(adjList)
    dist = [INF] * n
    dist[start] = 0
    queue = deque([start])
    pre = [(-1, -1)] * n  # (preNode, preEdge) bfs记录路径

    while queue:
        cur = queue.popleft()
        for next, edge in adjList[cur]:
            if edge == banEdge:
                continue
            cand = dist[cur] + 1
            if cand < dist[next]:
                pre[next] = (cur, edge)  # type: ignore
                dist[next] = cand
                queue.append(next)

    return dist, pre


n, m = map(int, input().split())
adjList = [[] for _ in range(n)]
for i in range(m):
    u, v = map(int, input().split())
    u, v = u - 1, v - 1
    adjList[u].append((v, i))  # !记录每条边的编号

dist, pre = bfs2(0, adjList, -1)
if dist[n - 1] == INF:
    print(-1)
    exit(0)

path = []  # !记录最短路上的边id
cur = n - 1
while cur != -1:
    path.append(pre[cur][1])
    cur = pre[cur][0]

path.reverse()
