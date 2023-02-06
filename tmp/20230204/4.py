from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的 m x n 二进制 矩阵 grid 。你可以从一个格子 (row, col) 移动到格子 (row + 1, col) 或者 (row, col + 1) ，前提是前往的格子值为 1 。如果从 (0, 0) 到 (m - 1, n - 1) 没有任何路径，我们称该矩阵是 不连通 的。

# 你可以翻转 最多一个 格子的值（也可以不翻转）。你 不能翻转 格子 (0, 0) 和 (m - 1, n - 1) 。

# 如果可以使矩阵不连通，请你返回 true ，否则返回 false 。

# 注意 ，翻转一个格子的值，可以使它的值从 0 变 1 ，或从 1 变 0 。

from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(key) for key in self.parent))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


def findCutVertices(n: int, graph: List[List[int]]) -> List[bool]:
    """Tarjan 算法求无向图的割点

    Args:
        n (int): 顶点数
        graph (List[List[int]]): 邻接表

    Returns:
        List[bool]: 每个点是否是割点
    """

    def dfs(cur: int, pre: int) -> int:
        nonlocal dfsId
        dfsId += 1
        dfsOrder[cur] = dfsId
        curLow = dfsId
        childCount = 0
        for next in graph[cur]:
            if dfsOrder[next] == 0:
                childCount += 1
                nextLow = dfs(next, cur)
                if nextLow >= dfsOrder[cur]:
                    isCut[cur] = True
                if nextLow < curLow:
                    curLow = nextLow
            elif next != pre and dfsOrder[next] < curLow:
                curLow = dfsOrder[next]
        if pre == -1 and childCount == 1:  # 特判：只有一个儿子的树根，删除后并没有增加连通分量的个数，这种情况下不是割顶
            isCut[cur] = False
        return curLow

    isCut = [False] * n
    dfsOrder = [0] * n  # 值从 1 开始
    dfsId = 0
    for i, order in enumerate(dfsOrder):
        if order == 0:
            dfs(i, -1)

    return isCut


DIR2 = ((0, 1), (1, 0))


class Solution:
    def isPossibleToCutPath(self, grid: List[List[int]]) -> bool:
        ROW, COL = len(grid), len(grid[0])
        uf = UnionFindArray(ROW * COL)
        adjList = [[] for _ in range(ROW * COL)]
        for r in range(ROW):
            for c in range(COL):
                if grid[r][c] == 0:
                    continue
                for dr, dc in DIR2:
                    nr, nc = r + dr, c + dc
                    if 0 <= nr < ROW and 0 <= nc < COL and grid[nr][nc] == 1:
                        uf.union(r * COL + c, nr * COL + nc)
                        adjList[r * COL + c].append(nr * COL + nc)
                        adjList[nr * COL + nc].append(r * COL + c)

        if not uf.isConnected(0, (ROW - 1) * COL + COL - 1):
            return True

        isCut = findCutVertices(ROW * COL, adjList)
        return any(isCut)
