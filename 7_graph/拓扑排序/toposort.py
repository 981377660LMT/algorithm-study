"""
最好是根据实际场景写逻辑 而不是复用模板
"""

from typing import DefaultDict, Hashable, List, Set, Tuple, TypeVar
from collections import defaultdict, deque


def topoSort(n: int, adjList: List[List[int]], deg: List[int], directed: bool) -> List[int]:
    """求图的拓扑排序

    Args:
        n (int): 顶点0~n-1
        adjList (List[List[int]]): 邻接表
        deg (List[int]): 有向图的入度/无向图的度
        directed (bool): 是否为有向图

    Returns:
        List[int]: 拓扑排序结果, 若不存在则返回空列表
    """
    startDeg = 0 if directed else 1
    queue = deque([v for v in range(n) if deg[v] == startDeg])
    res = []
    while queue:
        cur = queue.popleft()
        res.append(cur)
        for next in adjList[cur]:
            deg[next] -= 1
            if deg[next] == startDeg:
                queue.append(next)

    return [] if any(deg) else res
    return [] if len(res) < n else res


def toposort1(n: int, adjMap: DefaultDict[int, Set[int]], deg: List[int]) -> Tuple[int, List[int]]:
    """ "返回有向图拓扑排序方案数和拓扑排序结果

    注意图里有重边时不能多次计算deg
    """
    queue = deque([v for v in range(n) if deg[v] == 0])
    res = []
    topoCount = 1
    while queue:
        topoCount *= len(queue)
        cur = queue.popleft()
        res.append(cur)
        for next in adjMap[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)
    if len(res) != n:
        return 0, []
    return topoCount, res


T = TypeVar("T", bound=Hashable)


def toposort2(
    adjMap: DefaultDict[T, Set[T]], deg: DefaultDict[T, int], /, allVertex: Set[T]
) -> Tuple[int, List[T]]:
    """ "返回有向图拓扑排序方案数和拓扑排序结果

    注意图里有重边时不能多次计算deg
    """
    for v in allVertex:  # !初始化所有顶点的入度
        deg.setdefault(v, 0)
    queue = deque([v for v in allVertex if deg[v] == 0])
    res, topoCount = [], 1
    while queue:
        topoCount *= len(queue)
        cur = queue.popleft()
        res.append(cur)
        for next in adjMap[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)
    if len(res) != len(allVertex):
        return 0, []
    return topoCount, res


def toposort3(preDeps: List[Tuple[T, T]], allVertex: Set[T]) -> Tuple[int, List[int]]:
    """ "返回有向图拓扑排序方案数和拓扑排序结果

    注意图里有重边时不能多次计算deg
    """
    adjMap = defaultdict(set)
    deg = defaultdict(int)
    for pre, next in preDeps:
        # !保证重边时deg不重复计算
        if next not in adjMap[pre]:
            adjMap[pre].add(next)
            deg[next] += 1

    for v in allVertex:
        deg.setdefault(v, 0)

    queue = deque([v for v in allVertex if deg[v] == 0])
    res = []
    topoCount = 1
    while queue:
        topoCount *= len(queue)
        cur = queue.popleft()
        res.append(cur)
        for next in adjMap[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)

    if len(res) != len(allVertex):
        return 0, []
    return topoCount, res
