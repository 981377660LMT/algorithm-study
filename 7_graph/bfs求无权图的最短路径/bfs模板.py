# bfs模板 类似于dijkstra
from collections import defaultdict, deque
from typing import Hashable, TypeVar


INF = int(1e20)
T = TypeVar("T", bound=Hashable)


def bfs(adjMap: defaultdict[T, set[T]], start: T) -> defaultdict[T, int]:
    """时间复杂度O(V+E)"""
    dist = defaultdict(lambda: INF, {start: 0})
    queue: deque[tuple[int, T]] = deque([(0, start)])

    while queue:
        _, cur = queue.popleft()
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + 1:
                dist[next] = dist[cur] + 1
                queue.append((dist[next], next))

    return dist
