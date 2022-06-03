# 无向图联通分量:
# dfs
# 并查集

from collections import defaultdict
from typing import DefaultDict, List


class UnionFindArray:
    def __init__(self, n: int):
        self.n = n
        self.count = n
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
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups


class Solution:
    def countComponents2(self, n: int, edges: List[List[int]]) -> int:
        """并查集求无向图连通分量"""
        uf = UnionFindArray(n)
        for edge in edges:
            uf.union(edge[0], edge[1])
        return uf.count

    def countComponents(self, n: int, edges: List[List[int]]) -> int:
        """dfs求无向图连通分量"""

        def dfs(cur: int) -> None:
            if cur in visited:
                return
            visited.add(cur)
            for next in adjMap[cur]:
                dfs(next)

        adjMap = defaultdict(set)
        for u, v in edges:
            adjMap[u].add(v)
            adjMap[v].add(u)

        visited = set()
        res = 0
        for i in range(n):
            if i not in visited:
                dfs(i)
                res += 1
        return res
