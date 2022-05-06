# https://oi-wiki.org/graph/min-circle/#floyd
from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List


def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    """时间复杂度O((V+E)logV)"""
    dist = [int(1e20)] * n
    dist[start] = 0  # 注意这里不要忘记初始化pq里的
    pq = [(0, start)]

    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:  # 剪枝，有的题目不加就TLE
            continue
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))

    return int(1e20)


def minCycle(n: int, edges: List[List[int]]) -> int:
    """枚举所有边，每一次求删除一条边之后对这条边的起点跑一次 Dijkstra"""
    adjMap = defaultdict(lambda: defaultdict(lambda: int(1e20)))
    for u, v, w in edges:
        adjMap[u][v] = min(adjMap[u][v], w)
        adjMap[v][u] = min(adjMap[v][u], w)

    res = int(1e20)
    for u, v, _ in edges:
        tmp = adjMap[u][v]
        adjMap[u][v] = int(1e20)
        adjMap[v][u] = int(1e20)
        res = min(res, dijkstra(n, adjMap, u, v) + tmp)
        adjMap[u][v] = tmp
        adjMap[v][u] = tmp

    return res

