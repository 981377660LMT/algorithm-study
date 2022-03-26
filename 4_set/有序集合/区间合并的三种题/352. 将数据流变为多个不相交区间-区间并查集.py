from typing import List
from sortedcontainers import SortedSet


class SummaryRanges:
    def __init__(self):
        self.uf = UnionFindArray(int(1e4) + 10)
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

        rootX, rootY = sorted([rootX, rootY], reverse=True)
        # 小的总是指向大的
        self.parent[rootY] = rootX
        self.rank[rootX] += self.rank[rootY]
        self.count -= 1
        return True

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)
