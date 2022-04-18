"""dijkstra模板"""

from collections import defaultdict
from functools import lru_cache
from heapq import heappop, heappush
from typing import DefaultDict, Hashable, List, Optional, Tuple, TypeVar, Union, overload

INF = int(1e20)


# T = TypeVar('T', bound=Hashable)


@overload
def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int) -> List[int]:
    ...


@overload
def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    ...


def dijkstra(
    n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: Optional[int] = None
) -> Union[int, List[int]]:
    """时间复杂度O((V+E)logV)"""
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
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
        _, cur, curPath = heappop(pq)
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
            if cur == end:
                return curDist
            for next in adjMap[cur]:
                if dist[next] > dist[cur] + adjMap[cur][next]:
                    dist[next] = dist[cur] + adjMap[cur][next]
                    heappush(pq, (dist[next], next))

        return INF

    return inner(start, end)
