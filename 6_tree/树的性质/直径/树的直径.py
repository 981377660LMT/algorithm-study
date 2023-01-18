from typing import List, Set, Tuple, Union
from collections import deque


def calDiameter(adjList: List[List[Tuple[int, int]]]) -> Tuple[int, List[int]]:
    """求带权树的(直径长度, 直径路径)"""

    def dfs(start: int) -> Tuple[int, List[int]]:
        dist = [-1] * n
        dist[start] = 0
        stack = [start]
        while stack:
            cur = stack.pop()
            for next, weight in adjList[cur]:
                if dist[next] != -1:
                    continue
                dist[next] = dist[cur] + weight
                stack.append(next)
        endPoint = dist.index(max(dist))
        return endPoint, dist

    n = len(adjList)
    u, _ = dfs(0)
    v, dist = dfs(u)
    diameter = dist[v]
    path = [v]
    while u != v:
        for next, weight in adjList[v]:
            if dist[next] + weight == dist[v]:
                path.append(next)
                v = next
                break

    return diameter, path


Tree = Union[List[List[int]], List[Set[int]]]


def calDiameter1(adjList: "Tree") -> Tuple[int, Tuple[int, int]]:
    """bfs计算树的直径长度和直径两端点"""
    n = len(adjList)
    queue = deque([0])
    visited = [False] * n
    visited[0] = True
    last1 = 0  # 第一次BFS最后一个点
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last1 = queue.popleft()
            for next in adjList[last1]:
                if not visited[next]:
                    visited[next] = True
                    queue.append(next)

    queue = deque([last1])  # 第一次最后一个点作为第二次BFS的起点
    visited = [False] * n
    visited[last1] = True
    last2 = 0  # 第二次BFS最后一个点
    res = -1
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last2 = queue.popleft()
            for next in adjList[last2]:
                if not visited[next]:
                    visited[next] = True
                    queue.append(next)
        res += 1

    return res, tuple(sorted([last1, last2]))


def calDiameter2(adjList: "Tree") -> List[int]:
    """dfs计算树的直径的`路径`"""

    def dfs(cur: int, pre: int) -> None:
        parent[cur] = pre
        depth[cur] = depth[pre] + 1
        for next in adjList[cur]:
            if next != pre:
                dfs(next, cur)

    n = len(adjList)
    depth = [0] * n
    parent = [-1] * n

    dfs(0, -1)
    v1 = depth.index(max(depth))
    dfs(v1, -1)
    v2 = depth.index(max(depth))
    path = []
    while v2 != -1:
        path.append(v2)
        v2 = parent[v2]
    return path


def calDiameter3(adjList: "Tree") -> List[int]:
    """bfs计算树的直径的`路径`"""

    n = len(adjList)
    queue = deque([0])
    visited = [False] * n
    visited[0] = True
    last1 = 0
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last1 = queue.popleft()
            for next in adjList[last1]:
                if not visited[next]:
                    visited[next] = True
                    queue.append(next)

    queue = deque([(last1, -1)])
    visited = [False] * n
    visited[last1] = True
    last2 = 0
    depth = [0] * n
    parent = [-1] * n
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last2, pre = queue.popleft()
            parent[last2] = pre
            depth[last2] = depth[pre] + 1
            for next in adjList[last2]:
                if not visited[next]:
                    visited[next] = True
                    queue.append((next, last2))  # type: ignore

    path = []
    while last2 != -1:
        path.append(last2)
        last2 = parent[last2]
    return path


if __name__ == "__main__":
    edges = [[0, 1], [1, 2], [2, 3]]
    adjList: List[List[int]] = [[] for _ in range(4)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    assert calDiameter1(adjList) == (3, (0, 3))
    assert len(calDiameter2(adjList)) == 4
    assert len(calDiameter3(adjList)) == 4
