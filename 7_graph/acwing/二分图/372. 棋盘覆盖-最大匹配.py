# 给定一个 N 行 N 列的棋盘，已知某些格子禁止放置。
# 求最多能往棋盘上放多少块的长度为 2、宽度为 1 的骨牌，骨牌的边界与格线重合（骨牌占用两个格子），并且任意两张骨牌都不重叠。

# 第一行包含两个整数 N 和 t，其中 t 为禁止放置的格子的数量。
# 接下来 t 行每行包含两个整数 x 和 y，表示位于第 x 行第 y 列的格子禁止放置，行列数从 1 开始。
from collections import defaultdict
from typing import DefaultDict, List, Set
import sys


sys.setrecursionlimit(1000000)


def hungarian(adjMap: DefaultDict[int, Set[int]]) -> int:
    def getColor(adjMap: DefaultDict[int, Set[int]]) -> DefaultDict[int, int]:
        """检测二分图并染色"""

        def dfs(cur: int, color: int) -> None:
            colors[cur] = color
            for next in adjMap[cur]:
                if colors[next] == -1:
                    dfs(next, color ^ 1)
                elif colors[cur] == colors[next]:
                    raise Exception('不是二分图')

        colors = defaultdict(lambda: -1)
        for i in range(n):
            if colors[i] == -1:
                dfs(i, 0)
        return colors

    def dfs(boy: int) -> bool:
        """寻找增广路"""
        nonlocal visited
        if visited[boy]:
            return False
        visited[boy] = True

        for girl in adjMap[boy]:
            if matching[girl] == -1 or dfs(matching[girl]):
                matching[boy] = girl
                matching[girl] = boy
                return True
        return False

    n = len(adjMap)
    maxMatching = 0
    matching = [-1] * n
    colors = getColor(adjMap)
    visited = [False] * n
    for i in range(n):
        visited = [False] * n
        if colors[i] == 0 and matching[i] == -1:
            if dfs(i):
                maxMatching += 1

    return maxMatching


n, bad = map(int, input().split())
badSet = set()
for _ in range(bad):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    badSet.add(x * n + y)

adjList = [[] for _ in range(n * n)]
# build graph
for i in range(n):
    for j in range(n):
        cur = i * n + j
        if cur not in badSet:
            if i > 0:
                next = (i - 1) * n + j
                if next not in badSet:
                    adjList[cur].append(next)
                    adjList[next].append(cur)
            if j > 0:
                next = i * n + (j - 1)
                if next not in badSet:
                    adjList[cur].append(next)
                    adjList[next].append(cur)

# 匈牙利算法
print(hungarian(adjList))

