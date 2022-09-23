from graphlib import CycleError, TopologicalSorter

from typing import List
from collections import deque


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


class Solution:
    def canFinish(self, numCourses: int, prerequisites: list[list[int]]) -> bool:
        """有向图是否无环"""
        ts = TopologicalSorter()
        for cur, pre in prerequisites:
            ts.add(cur, pre)
        try:
            ts.prepare()
            return True
        except CycleError:
            return False

    def canFinish2(self, numCourses: int, prerequisites: list[list[int]]) -> bool:
        """有向图是否无环"""
        adjList = [[] for _ in range(numCourses)]
        deg = [0] * numCourses
        for cur, pre in prerequisites:
            adjList[pre].append(cur)
            deg[cur] += 1
        return len(topoSort(numCourses, adjList, deg, True)) == numCourses


print(Solution().canFinish(2, [[1, 0], [0, 1]]))
print(Solution().canFinish(2, [[1, 0]]))
print(Solution().canFinish2(2, [[1, 0], [0, 1]]))
print(Solution().canFinish2(2, [[1, 0]]))
