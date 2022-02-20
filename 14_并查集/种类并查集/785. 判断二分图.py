from typing import List


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


class Solution:
    def isBipartite(self, graph: List[List[int]]) -> bool:
        n = len(graph)
        uf = UnionFindArray(n * 2)
        for cur, nexts in enumerate(graph):
            for next in nexts:
                if uf.isConnected(cur, next):
                    return False
                uf.union(cur, next + n)
                uf.union(cur + n, next)
        return True
