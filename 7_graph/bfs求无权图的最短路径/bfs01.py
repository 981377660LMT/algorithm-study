# bfs01/01bfs


from collections import deque
from typing import List, Tuple

INF = int(1e18)


def bfs01(adjList: List[List[Tuple[int, int]]], start: int) -> Tuple[List[int], List[int]]:
    """返回 (最短路径长度,每个点的前驱节点)"""
    n = len(adjList)
    dist = [INF] * n
    pre = [-1] * n

    queue = deque([start])
    dist[start] = 0
    while queue:
        cur = queue.popleft()
        for next, cost in adjList[cur]:
            cand = dist[cur] + cost
            if cand < dist[next]:
                dist[next] = cand
                pre[next] = cur
                if cost == 0:
                    queue.appendleft(next)
                else:
                    queue.append(next)
    return dist, pre


def bfs01MultiStart(
    adjList: List[List[Tuple[int, int]]], starts: List[int]
) -> Tuple[List[int], List[int], List[int]]:
    """多个起点的01bfs, 返回 (最短路径长度,每个点的前驱节点,每个点的起点)"""
    n = len(adjList)
    dist = [INF] * n
    pre = [-1] * n
    root = [-1] * n
    queue = deque(starts)
    for start in starts:
        dist[start] = 0
        root[start] = start

    while queue:
        cur = queue.popleft()
        for next, cost in adjList[cur]:
            cand = dist[cur] + cost
            if cand < dist[next]:
                dist[next] = cand
                root[next] = root[cur]
                pre[next] = cur
                if cost == 0:
                    queue.appendleft(next)
                else:
                    queue.append(next)

    return dist, pre, root
