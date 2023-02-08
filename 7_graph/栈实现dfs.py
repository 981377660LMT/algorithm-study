"""
栈实现dfs
适用于atcoder
"""
# dfs递归改成栈实现
# https://tjkendev.github.io/procon-library/python/graph/dfs.html
from typing import List


n = 5
edges = [[0, 1], [0, 2], [1, 3], [1, 4]]
adjList = [[] for _ in range(n)]
for u, v in edges:
    adjList[u].append(v)
    adjList[v].append(u)


def stackDfs(adjList: List[List[int]], start: int, target: int) -> List[int]:
    """栈实现dfs,求出start到target的路径"""
    n = len(adjList)
    stack = [(start, -1, 0)]  # cur, pre, dep
    path = [0] * n  # !记录路径(每个深度对应的结点)
    while stack:
        cur, pre, dep = stack.pop()
        path[dep] = cur

        # !处理当前结点的逻辑
        if cur == target:
            return path[: dep + 1]

        for next in adjList[cur]:
            if next == pre:
                continue
            stack.append((next, cur, dep + 1))

    return []
