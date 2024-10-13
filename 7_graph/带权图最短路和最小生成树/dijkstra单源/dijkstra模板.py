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


# !dijkstra求出路径上每个点的前驱点和前驱边
# 其中邻接表的每个元素是一个三元组，分别是邻接点，边权，边的编号
def dijkstra3(
    n: int, adjList: Sequence[Sequence[Tuple[int, int, int]]], start: int
) -> Tuple[List[int], List[int], List[int]]:
    dist = [INF] * n
    preV, preE = [-1] * n, [-1] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight, eid in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                preV[next] = cur
                preE[next] = eid
                heappush(pq, (dist[next], next))
    return dist, preV, preE


# 多源最短路, 返回(距离, 前驱, 根节点).用于求出离每个点最近的起点.
def dijkstraMultiRoot(
    n: int, adjList: List[List[Tuple[int, int]]], roots: List[int]
) -> Tuple[List[int], List[int], List[int]]:
    dist = [INF] * n
    pre = [-1] * n
    root = [-1] * n
    pq = [(0, v) for v in roots]
    for v in roots:
        dist[v] = 0
        root[v] = v
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                root[next] = root[cur]
                pre[next] = cur
                heappush(pq, (dist[next], next))
    return dist, pre, root


def dijkstraWithCount(
    n: int, adjList: List[List[Tuple[int, int]]], start: int, *, mod=int(1e9 + 7)
) -> Tuple[List[int], List[int]]:
    """dijkstra求出起点到各点的最短距离和最短路径数量(模mod).

    时间复杂度O((V+E)logV).
    """
    dist = [INF] * n
    count = [0] * n
    dist[start] = 0
    count[start] = 1
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next, weight in adjList[cur]:
            cand = dist[cur] + weight
            if cand < dist[next]:
                dist[next] = cand
                count[next] = count[cur]
                heappush(pq, (dist[next], next))
            elif cand == dist[next]:
                count[next] += count[cur]
                if count[next] >= mod:
                    count[next] -= mod
    return dist, count


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
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
