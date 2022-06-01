from typing import DefaultDict, Hashable, List, Set, Tuple, TypeVar
from collections import defaultdict, deque


def toposort1(n: int, adjMap: DefaultDict[int, Set[int]], deg: List[int]) -> Tuple[int, List[int]]:
    """"返回有向图拓扑排序方案数和拓扑排序结果
    
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


T = TypeVar('T', bound=Hashable)


def toposort2(
    adjMap: DefaultDict[T, Set[T]], deg: DefaultDict[T, int], /, allVertex: Set[T]
) -> Tuple[int, List[T]]:
    """"返回有向图拓扑排序方案数和拓扑排序结果
    
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
    """"返回有向图拓扑排序方案数和拓扑排序结果
    
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
