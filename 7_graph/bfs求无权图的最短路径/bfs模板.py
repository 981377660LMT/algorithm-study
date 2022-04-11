# bfs模板 类似于dijkstra
from collections import defaultdict, deque
from typing import DefaultDict, Deque, Optional, Set, Tuple, Union, overload

MOD = int(1e9 + 7)
INF = int(1e20)


@overload
def bfs(adjMap: DefaultDict[int, Set[int]], start: int) -> DefaultDict[int, int]:
    ...


@overload
def bfs(adjMap: DefaultDict[int, Set[int]], start: int, end: int) -> int:
    ...


def bfs(
    adjMap: DefaultDict[int, Set[int]], start: int, end: Optional[int] = None
) -> Union[int, DefaultDict[int, int]]:
    """时间复杂度O(V+E)"""
    dist = defaultdict(lambda: INF, {key: INF for key in adjMap.keys()})
    dist[start] = 0
    queue: Deque[Tuple[int, int]] = deque([(0, start)])

    while queue:
        curDist, cur = queue.popleft()
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + 1:
                dist[next] = dist[cur] + 1
                queue.append((dist[next], next))

    return INF if end is not None else dist
