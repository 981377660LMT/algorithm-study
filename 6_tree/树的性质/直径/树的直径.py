from typing import List, Set, Tuple, Union
from collections import deque


def calDiameter(adjList: List[List[Tuple[int, int]]], start=0) -> Tuple[int, List[int]]:
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
    u, _ = dfs(start)
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


def getTreeDiameter(n: int, tree: List[List[int]], start=0) -> Tuple[int, List[int]]:
    """求无权树的(直径长度,直径路径)."""

    def dfs(start: int) -> Tuple[int, List[int]]:
        dist = [-1] * n
        dist[start] = 0
        stack = [start]
        while stack:
            cur = stack.pop()
            for next in tree[cur]:
                if dist[next] != -1:
                    continue
                dist[next] = dist[cur] + 1
                stack.append(next)
        endPoint = dist.index(max(dist))
        return endPoint, dist

    u, _ = dfs(start)
    v, dist = dfs(u)
    diameter = dist[v]
    path = [v]
    while u != v:
        for next in tree[v]:
            if dist[next] + 1 == dist[v]:
                path.append(next)
                v = next
                break

    return diameter, path


Tree = Union[List[List[int]], List[Set[int]]]


def calDiameter1(adjList: "Tree", start=0) -> Tuple[int, Tuple[int, int]]:
    """bfs计算树的直径长度和直径两端点"""
    n = len(adjList)
    queue = deque([start])
    visited = [False] * n
    visited[start] = True
    last1 = start  # 第一次BFS最后一个点
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


def calDiameter2(adjList: "Tree", start=0) -> List[int]:
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

    dfs(start, -1)
    v1 = depth.index(max(depth))
    dfs(v1, -1)
    v2 = depth.index(max(depth))
    path = []
    while v2 != -1:
        path.append(v2)
        v2 = parent[v2]
    return path


def calDiameter3(adjList: "Tree", start=0) -> List[int]:
    """bfs计算树的直径的`路径`"""

    n = len(adjList)
    queue = deque([start])
    visited = [False] * n
    visited[start] = True
    last1 = start
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

    # 100318. 合并两棵树后的最小直径(连边后树的最小直径) -> 连接直径中点
    # https://leetcode.cn/problems/find-minimum-diameter-after-merging-two-trees/description/
    class Solution:
        def minimumDiameterAfterMerge(
            self, edges1: List[List[int]], edges2: List[List[int]]
        ) -> int:
            n1, n2 = len(edges1) + 1, len(edges2) + 1
            adjList1 = [[] for _ in range(n1)]
            adjList2 = [[] for _ in range(n2)]
            for u, v in edges1:
                adjList1[u].append(v)
                adjList1[v].append(u)
            for u, v in edges2:
                adjList2[u].append(v)
                adjList2[v].append(u)

            d1, _ = getTreeDiameter(n1, adjList1)
            d2, _ = getTreeDiameter(n2, adjList2)

            return max(d1, d2, (d1 + 1) // 2 + (d2 + 1) // 2 + 1)
