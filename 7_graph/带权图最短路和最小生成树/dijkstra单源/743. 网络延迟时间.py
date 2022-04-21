from collections import defaultdict
from heapq import heappop, heappush
from typing import DefaultDict, List, Optional, Tuple, Union, overload

INF = int(1e20)


@overload
def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int) -> List[int]:
    ...


@overload
def dijkstra(n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: int) -> int:
    ...


def dijkstra(
    n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], start: int, end: Optional[int] = None
) -> Union[int, List[int]]:
    dist = [INF] * n
    dist[start] = 0
    pq = [(0, start)]
    while pq:
        curDist, cur = heappop(pq)
        if dist[cur] < curDist:
            continue
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + adjMap[cur][next]:
                dist[next] = dist[cur] + adjMap[cur][next]
                heappush(pq, (dist[next], next))
    return dist


class Solution:
    def networkDelayTime(self, times: List[List[int]], n: int, k: int) -> int:
        adjMap = defaultdict(lambda: defaultdict(lambda: INF))
        for u, v, w in times:
            adjMap[u - 1][v - 1] = w
        dist = dijkstra(n, adjMap, k - 1)
        return res if (res := max(dist)) < INF else -1

