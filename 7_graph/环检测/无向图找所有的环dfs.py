from typing import List
from typing import Iterable, Mapping, Sequence, Union

AdjList = Sequence[Iterable[int]]
AdjMap = Mapping[int, Iterable[int]]
Graph = Union[AdjList, AdjMap]


def findCycleGroup(n: int, graph: "Graph") -> List[List[int]]:
    """无向图找到所有的环，图可能不连通"""

    def dfs(cur: int, pre: int) -> bool:
        """环检测，并记录路径"""
        if visited[cur]:
            return True
        visited[cur] = True
        for next in graph[cur]:
            if next == pre:
                continue
            path.append(next)
            if dfs(next, cur):
                return True
            path.pop()
        return False

    def addCycleOnPath() -> None:
        curCycle = []
        last = path.pop()
        curCycle.append(last)
        while path and path[-1] != last:
            curCycle.append(path.pop())
        res.append(curCycle)

    res = []
    path = []
    visited = [False] * n

    for i in range(n):
        if visited[i]:
            continue
        path = [i]
        if dfs(i, -1):
            addCycleOnPath()

    return res
