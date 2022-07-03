from collections import deque
from typing import DefaultDict, List, Set, Tuple


def findCycleAndCalDepth1(
    n: int, adjMap: DefaultDict[int, DefaultDict[int, int]], indeg: List[int]
) -> Tuple[List[List[int]], List[int]]:
    """内向基环树找环上的点，并记录每个点在拓扑排序中的最大距离(带权)，最外层的点深度为0"""
    depth = [0] * n
    queue = deque([(i, 0) for i in range(n) if indeg[i] == 0])
    visited = [False] * n
    while queue:
        cur, dist = queue.popleft()
        visited[cur] = True
        for next in adjMap[cur]:
            depth[next] = max(depth[next], dist + adjMap[cur][next])
            indeg[next] -= 1
            if indeg[next] == 0:
                queue.append((next, dist + adjMap[cur][next]))

    def dfs(cur: int, path: List[int]) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        for next in adjMap[cur]:
            dfs(next, path)

    cycleGroup = []
    for i in range(n):
        if visited[i]:
            continue
        path = []
        dfs(i, path)
        cycleGroup.append(path)

    return cycleGroup, depth


def findCycleAndCalDepth2(
    n: int, adjMap: DefaultDict[int, Set[int]], deg: List[int]
) -> Tuple[List[List[int]], List[int]]:
    """无向基环树找环上的点，并记录每个点在拓扑排序中的最大深度，最外层的点深度为0"""
    depth = [0] * n
    queue = deque([(i, 0) for i in range(n) if deg[i] == 1])
    visited = [False] * n
    while queue:
        cur, dist = queue.popleft()
        visited[cur] = True
        for next in adjMap[cur]:
            depth[next] = max(depth[next], dist + 1)
            deg[next] -= 1
            if deg[next] == 1:
                queue.append((next, dist + 1))

    def dfs(cur: int, path: List[int]) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        path.append(cur)
        for next in adjMap[cur]:
            dfs(next, path)

    cycleGroup = []
    for i in range(n):
        if visited[i]:
            continue
        path = []
        dfs(i, path)
        cycleGroup.append(path)

    return cycleGroup, depth
