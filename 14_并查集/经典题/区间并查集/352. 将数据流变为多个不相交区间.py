from typing import List
from sortedcontainers import SortedSet

# 352. 将数据流变为多个不相交区间


class SummaryRanges:
    def __init__(self):
        self.uf = UnionFind(int(1e4) + 10)
        self.points = SortedSet()

    def addNum(self, val: int) -> None:
        self.uf.union(val, val + 1)
        self.points.add(val)

    def getIntervals(self) -> List[List[int]]:
        res = []
        for p in self.points:
            if res and p <= res[-1][1]:
                continue
            res.append([p, self.uf.find(p) - 1])
        return res


class UnionFind:
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
        """union后x所在的root的parent指向y所在的root"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False

        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.part -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)
