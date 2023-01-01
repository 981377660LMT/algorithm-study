from collections import deque
from typing import List
from typing import Iterable, Mapping, Sequence, Union

AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Graph = Union[AdjList, AdjMap]


def findCycle(n: int, graph: "Graph", degrees: List[int]) -> List[int]:
    """无向图找环上的点 拓扑排序，剪掉所有树枝"""
    queue = deque([i for i in range(n) if degrees[i] == 1])
    visited = [False] * n
    while queue:
        cur = queue.popleft()
        visited[cur] = True
        for next in graph[cur]:
            if visited[next]:
                continue
            degrees[next] -= 1
            if degrees[next] == 1:
                queue.append(next)

    onCycle = [i for i in range(n) if not visited[i]]
    return onCycle
