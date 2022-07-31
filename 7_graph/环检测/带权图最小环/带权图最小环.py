# https://oi-wiki.org/graph/min-circle/#floyd
"""O(V*E*Log(V+E))"""
from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List

INF = int(1e20)


def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    """时间复杂度O((V+E)logV)"""
    dist = [INF] * n
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

    return INF


def minCycle(n: int, edges: List[List[int]]) -> int:
    """枚举所有边，删除这条边之后以这条边的端点为起点终点跑一次 Dijkstra"""
    adjMap = defaultdict(lambda: defaultdict(lambda: INF))
    for u, v, w in edges:
        adjMap[u][v] = min(adjMap[u][v], w)
        adjMap[v][u] = min(adjMap[v][u], w)

    res = INF
    for u, v, _ in edges:
        tmp = adjMap[u][v]
        adjMap[u][v] = INF
        adjMap[v][u] = INF
        res = min(res, dijkstra(n, adjMap, u, v) + tmp)
        adjMap[u][v] = tmp
        adjMap[v][u] = tmp

    return res
