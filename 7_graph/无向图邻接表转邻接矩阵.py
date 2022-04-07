# 无向图邻接表转邻接矩阵
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


# 1 ≤ n * m ≤ 100,000
class Solution:
    def solve(self, adjList: List[List[int]]) -> List[List[int]]:
        n = len(adjList)
        uf = UnionFindArray(n)
        adjMatrix = [[0] * n for _ in range(n)]
        for cur, nexts in enumerate(adjList):
            for next in nexts:
                uf.union(cur, next)

        for i in range(n):
            for j in range(n):
                if uf.isConnected(i, j):
                    adjMatrix[i][j] = 1

        return adjMatrix

