"""dijk模板"""

from collections import defaultdict
from functools import lru_cache
from heapq import heappop, heappush
from typing import DefaultDict, Hashable, List, Tuple, TypeVar, overload

INF = int(1e20)
Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]


@overload
def dijkstra(adjMap: Graph, start: Vertex) -> DefaultDict[Vertex, int]:
    ...


@overload
def dijkstra(adjMap: Graph, start: Vertex, end: Vertex) -> int:
    ...


def dijkstra(adjMap: Graph, start: Vertex, end: Vertex | None = None):
    """时间复杂度O((V+E)logV)"""
    dist = defaultdict(lambda: INF)
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

    return INF if end is not None else dist


def dijkstra2(
    n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int
) -> Tuple[List[int], List[List[int]]]:
    """记录路径的dijk"""
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start, [start])]
    path = [[] for _ in range(n)]
    path[start] = [start]
    while pq:
        curDist, cur, curPath = heappop(pq)
        if dist[cur] < curDist:
            continue
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                path[next] = curPath + [next]
                heappush(pq, (dist[next], next, path[next]))
    return dist, path


# 字符串顶点的dijk
def dijkstra3(adjMap: DefaultDict[str, DefaultDict[str, int]], start: str, end: str) -> int:
    """时间复杂度O((V+E)logV)"""

    @lru_cache(None)
    def inner(start: str, end: str) -> int:
        pq = [(0, start)]
        dist = defaultdict(lambda: INF)
        dist[start] = 0

        while pq:
            curDist, cur = heappop(pq)
            if dist[cur] < curDist:
                continue
            if cur == end:
                return curDist
            for next in adjMap[cur]:
                if dist[next] > dist[cur] + adjMap[cur][next]:
                    dist[next] = dist[cur] + adjMap[cur][next]
                    heappush(pq, (dist[next], next))

        return INF

    return inner(start, end)
