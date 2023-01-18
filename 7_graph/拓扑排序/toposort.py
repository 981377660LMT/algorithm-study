"""
最好是根据实际场景写逻辑 而不是复用模板
"""

from typing import DefaultDict, Hashable, List, Mapping, MutableMapping, Set, Tuple, TypeVar
from collections import deque


def topoSort(n: int, adjList: List[List[int]], deg: List[int], directed=True) -> List[int]:
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

    return [] if len(res) < n else res


T = TypeVar("T", bound=Hashable)


def topoSort2(
    allVertex: Set[T],
    adjMap: Mapping[T, List[T]],
    deg: MutableMapping[T, int],
    directed=True,
) -> List[T]:
    """返回图的拓扑排序结果, 若不存在则返回空列表"""
    startDeg = 0 if directed else 1
    queue = deque([v for v in allVertex if deg[v] == startDeg])
    res = []
    while queue:
        cur = queue.popleft()
        res.append(cur)
        for next in adjMap[cur]:
            deg[next] -= 1
            if deg[next] == startDeg:
                queue.append(next)
    return [] if len(res) < len(allVertex) else res


def toposort3(
    allVertex: Set[T], adjMap: Mapping[T, List[T]], deg: DefaultDict[T, int], directed=True
) -> Tuple[int, List[T]]:
    """返回有向图拓扑排序方案数和拓扑排序结果"""
    startDeg = 0 if directed else 1
    queue = deque([v for v in allVertex if deg[v] == startDeg])
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
