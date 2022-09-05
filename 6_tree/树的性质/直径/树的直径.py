from heapq import nlargest
from typing import List, Tuple
from collections import deque


def calDiameter1(adjList: List[List[int]], start: int) -> Tuple[int, Tuple[int, int]]:
    """bfs计算树的直径长度和直径两端点"""
    queue = deque([start])
    visited = set([start])
    last1 = 0  # 第一次BFS最后一个点
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last1 = queue.popleft()
            for next in adjList[last1]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)

    queue = deque([last1])  # 第一次最后一个点作为第二次BFS的起点
    visited = set([last1])
    last2 = 0  # 第二次BFS最后一个点
    res = -1
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            last2 = queue.popleft()
            for next in adjList[last2]:
                if next not in visited:
                    visited.add(next)
                    queue.append(next)
        res += 1

    return res, tuple(sorted([last1, last2]))


def calDiameter2(adjList: List[List[int]], start: int) -> int:
    """dfs计算树的直径长度"""

    def dfs(cur: int, pre: int) -> int:
        nonlocal res
        cands = [0, 0]
        for next in adjList[cur]:
            if next != pre:
                cands.append(dfs(next, cur))
        max1, max2 = nlargest(2, cands)
        res = max(res, max1 + max2)
        return max(max1, max2) + 1

    res = 0
    dfs(start, -1)
    return res


if __name__ == "__main__":
    edges = [[0, 1], [1, 2], [2, 3]]
    adjList = [[] for _ in range(4)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)
    assert calDiameter1(adjList, 0) == (3, (0, 3))
    assert calDiameter2(adjList, 0) == 3
