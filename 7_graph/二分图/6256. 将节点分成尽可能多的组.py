"""
6256. 将节点分成尽可能多的组

给定一个无向图,可能是不连通的(暗示需要找到所有的连通分量)
请你将图划分为 m 个组（编号从 1 开始），满足以下要求：

1. 图中每个节点都只属于一个组。
2. 图中每条边连接的两个点 [ai, bi] ，如果 ai 属于编号为 x 的组,bi 属于编号为 y 的组，那么 |y - x| = 1 。
请你返回最多可以将节点分为多少个组（也就是最大的 m ）。如果没办法在给定条件下分组，请你返回 -1 。

解法:
!1. |y - x| = 1 <=> bfs在同一层的节点编号相同 <=> 不存在奇环 <=> 二分图检测
!3. 求每个连通分量的直径 <=> 暴力枚举连通分量中的每个起点做bfs,看最大层数

不存在奇环 <=> 二分图
求一般图的连通分量的直径 => 暴力枚举连通分量中的每个起点做bfs, 更新最大层数
"""

from collections import deque, defaultdict
from typing import DefaultDict, List


class Solution:
    def magnificentSets(self, n: int, edges: List[List[int]]) -> int:
        adjList = [[] for _ in range(n)]
        uf = UnionFind(n)
        for u, v in edges:
            u, v = u - 1, v - 1
            adjList[u].append(v)
            adjList[v].append(u)
            uf.union(u, v)
        if not isBipartite(n, adjList):
            return -1
        return sum(calDiameter(n, adjList, group) + 1 for group in uf.getGroups().values())


def calDiameter(n: int, adjList: List[List[int]], group: List[int]) -> int:
    """bfs求连通分量 `group` 的直径长度"""
    res = 0
    for start in group:
        visited, queue = set([start]), deque([start])
        diameter = -1
        while queue:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                for next in adjList[cur]:
                    if next in visited:
                        continue
                    visited.add(next)
                    queue.append(next)
            diameter += 1
        res = max(res, diameter)
    return res


def isBipartite(n: int, adjList: List[List[int]]) -> bool:
    """二分图检测 dfs染色"""

    def dfs(cur: int, color: int) -> bool:
        colors[cur] = color
        for next in adjList[cur]:
            if colors[next] == -1:
                if not dfs(next, color ^ 1):
                    return False
            elif colors[next] == color:
                return False
        return True

    colors = [-1] * n
    for i in range(n):
        if colors[i] == -1 and not dfs(i, 0):
            return False
    return True


class UnionFind:
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

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


assert Solution().magnificentSets(n=6, edges=[[1, 2], [1, 4], [1, 5], [2, 6], [2, 3], [4, 6]]) == 4
