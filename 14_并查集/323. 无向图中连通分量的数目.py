# 无向图联通分量:
# dfs
# 并查集

from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

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


class Solution:
    def countComponents2(self, n: int, edges: List[List[int]]) -> int:
        """并查集求无向图连通分量"""
        uf = UnionFindArray(n)
        for edge in edges:
            uf.union(edge[0], edge[1])
        return uf.part

    def countComponents(self, n: int, edges: List[List[int]]) -> List[List[int]]:
        """
        dfs求无向图连通分量
        注意:这样求出来的group的相邻元素是直接相连的,与并查集求连通分量不同
        """

        def dfs(cur: int) -> None:
            if visited[cur]:
                return
            visited[cur] = True
            group.append(cur)
            for next in adjList[cur]:
                dfs(next)

        adjList = [[] for _ in range(n)]
        for u, v in edges:
            adjList[u].append(v)
            adjList[v].append(u)

        visited = [False] * n
        res = []
        for i in range(n):
            if not visited[i]:
                group = []
                dfs(i)
                res.append(group)
        return res
