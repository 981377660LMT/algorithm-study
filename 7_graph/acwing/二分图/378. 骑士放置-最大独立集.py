# 给定一个 N×M 的棋盘，有一些格子禁止放棋子。
# 问棋盘上最多能放多少个不能互相攻击的骑士（国际象棋的“骑士”，类似于中国象棋的“马”，按照“日”字攻击，但没有中国象棋“别马腿”的规则）。

# 第一行包含三个整数 N,M,T，其中 T 表示禁止放置的格子的数量。
# 接下来 T 行每行包含两个整数 x 和 y，表示位于第 x 行第 y 列的格子禁止放置，行列数从 1 开始。
# 输入格式

from collections import defaultdict
from typing import DefaultDict, List, Set
import sys


sys.setrecursionlimit(1000000)


def hungarian(adjMap: DefaultDict[int, Set[int]]) -> int:
    def getColor(adjMap: DefaultDict[int, Set[int]]) -> List[int]:
        """检测二分图并染色"""

        def dfs(cur: int, color: int) -> None:
            colors[cur] = color
            for next in adjMap[cur]:
                if colors[next] == -1:
                    dfs(next, color ^ 1)
                elif colors[cur] == colors[next]:
                    raise Exception('不是二分图')

        n = len(adjMap)
        colors = [-1] * n
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


if __name__ == '__main__':
    n, m, t = map(int, input().split())
    badSet = set()
    for _ in range(t):
        x, y = map(int, input().split())
        x, y = x - 1, y - 1
        badSet.add(x * n + y)

    dirs = [(1, 2), (1, -2), (-1, 2), (-1, -2), (2, 1), (2, -1), (-2, 1), (-2, -1)]
    adjList = defaultdict(set)
    # build graph
    for i in range(n):
        for j in range(m):
            cur = i * m + j
            if cur not in badSet:
                for dx, dy in dirs:
                    nextX, nextY = i + dx, j + dy
                    if 0 <= nextX < n and 0 <= nextY < m:
                        next = nextX * m + nextY
                        if next not in badSet:
                            adjList[cur].add(next)
                            adjList[next].add(cur)

    # 匈牙利算法
    print(n * m - hungarian(adjList))
