"""dijk模板"""

from collections import defaultdict
from functools import lru_cache
from heapq import heappop, heappush
from typing import DefaultDict, Hashable, List, Optional, Tuple, TypeVar, overload

INF = int(1e20)
Vertex = TypeVar('Vertex', bound=Hashable)
Graph = DefaultDict[Vertex, DefaultDict[Vertex, int]]


@overload
def dijkstra(adjMap: Graph, start: Vertex) -> DefaultDict[Vertex, int]:
    ...


@overload
def dijkstra(adjMap: Graph, start: Vertex, end: Vertex) -> int:
    ...


def dijkstra(adjMap: Graph, start: Vertex, end: Optional[Vertex] = None):
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


##########################################################################
def dijkstra2(
    n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int
) -> List[int]:
    """记录路径的dijk 用pre数组记录路径"""
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    pre = [start] * n
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        if cur == end:
            break
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                pre[next] = cur
                heappush(pq, (dist[next], next))

    res = [end]
    cur = end
    while cur != start:
        cur = pre[cur]
        res.append(cur)

    return res[::-1]


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
