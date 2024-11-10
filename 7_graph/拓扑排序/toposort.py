from typing import List, Tuple
from collections import deque


def hasCycle(n: int, adjList: List[List[int]], directed=True) -> bool:
    """拓扑排序判环."""
    if directed:
        deg = [0] * n
        for i in range(n):
            for j in adjList[i]:
                deg[j] += 1
    else:
        deg = [len(adj) for adj in adjList]

    startDeg = 0 if directed else 1
    queue = deque([v for v in range(n) if deg[v] == startDeg])
    count = 0
    while queue:
        cur = queue.popleft()
        count += 1
        for next in adjList[cur]:
            deg[next] -= 1
            if deg[next] == startDeg:
                queue.append(next)
    return count < n


def topoSort(n: int, adjList: List[List[int]], directed=True) -> Tuple[List[int], bool]:
    """求图的拓扑排序."""
    if directed:
        deg = [0] * n
        for i in range(n):
            for j in adjList[i]:
                deg[j] += 1
    else:
        deg = [len(adj) for adj in adjList]

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

    if len(res) != n:
        return [], False
    return res, True


from heapq import heapify, heappop, heappush
from typing import List, Tuple


def topoSortByHeap(
    n: int, adjList: List[List[int]], directed=True, minFirst=True
) -> Tuple[List[int], bool]:
    """使用优先队列的拓扑排序."""
    if directed:
        deg = [0] * n
        for i in range(n):
            for j in adjList[i]:
                deg[j] += 1
    else:
        deg = [len(adj) for adj in adjList]

    startDeg = 0 if directed else 1
    pq = [v if minFirst else -v for v in range(n) if deg[v] == startDeg]
    heapify(pq)
    res = []

    if minFirst:
        while pq:
            cur = heappop(pq)
            res.append(cur)
            for next in adjList[cur]:
                deg[next] -= 1
                if deg[next] == startDeg:
                    heappush(pq, next)
    else:
        while pq:
            cur = -heappop(pq)
            res.append(cur)
            for next in adjList[cur]:
                deg[next] -= 1
                if deg[next] == startDeg:
                    heappush(pq, -next)
    if len(res) != n:
        return [], False
    return res, True


if __name__ == "__main__":

    class Solution:
        def canFinish(self, numCourses: int, prerequisites: List[List[int]]) -> bool:
            n = numCourses
            adjList = [[] for _ in range(n)]
            for pre, cur in prerequisites:
                adjList[pre].append(cur)
            return not hasCycle(n, adjList)


# T = TypeVar("T", bound=Hashable)


# def toposort3(
#     allVertex: Set[T], adjMap: Mapping[T, List[T]], deg: DefaultDict[T, int], directed=True
# ) -> Tuple[int, List[T]]:
#     """返回有向图拓扑排序方案数和拓扑排序结果"""
#     startDeg = 0 if directed else 1
#     queue = deque([v for v in allVertex if deg[v] == startDeg])
#     visited = set()
#     res, topoCount = [], 1
#     while queue:
#         topoCount *= len(queue)
#         cur = queue.popleft()
#         res.append(cur)
#         visited.add(cur)
#         for next in adjMap[cur]:
#             if next in visited:
#                 continue
#             deg[next] -= 1
#             if deg[next] == 0:
#                 queue.append(next)
#     if len(res) != len(allVertex):
#         return 0, []
#     return topoCount, res
