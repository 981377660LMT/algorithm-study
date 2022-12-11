"""
三种颜色给图染色的方案数
给定一个无向图 n 个点 m 条边 (n<=20)
现在要给每个点染色 使得相邻的点颜色不同

!在每个连通分量里染色 dfs枚举 时间复杂度O(n^2*2^n)
"""

import sys
from typing import List, Tuple

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")


def rgbColoring2(n: int, edges: List[Tuple[int, int]]) -> int:
    """三种颜色给图染色 相邻点颜色不同的方案数"""
    adjList = [[] for _ in range(n)]
    for u, v in edges:
        adjList[u].append(v)
        adjList[v].append(u)

    res = 1
    groups = getGroups(n, adjList)
    for group in groups:
        res *= cal(adjList, group)
    return res


def cal(adjList: List[List[int]], group: List[int]) -> int:
    """连通分量  `group` 里的染色方案数"""

    def dfs(index: int) -> None:
        nonlocal res
        if index == len(group):
            res += 1
            return
        cur = group[index]
        for color in range(3):
            colors[cur] = color
            if all(colors[next] != color for next in adjList[cur]):
                dfs(index + 1)
            colors[cur] = -1

    res = 0
    colors = [-1] * n
    dfs(0)
    return res


def getGroups(n: int, adjList: List[List[int]]) -> List[List[int]]:
    """dfs求连通分量"""

    def dfs(cur: int) -> None:
        if visited[cur]:
            return
        visited[cur] = True
        group.append(cur)
        for next in adjList[cur]:
            dfs(next)

    visited = [False] * n
    res = []
    for i in range(n):
        if not visited[i]:
            group = []
            dfs(i)
            res.append(group)
    return res


if __name__ == "__main__":

    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        a, b = map(int, input().split())
        edges.append((a - 1, b - 1))

    res = rgbColoring2(n, edges)
    print(res)
