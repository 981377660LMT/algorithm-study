# 敌人开2*n即可
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


# 互相不喜欢的人可以用特殊的并查集，每个集分为 2 部分，不喜欢的人一定在同一个集的另一部分里面
class Solution:
    def possibleBipartition(self, n: int, dislikes: List[List[int]]) -> bool:
        uf = UnionFindArray(n * 2 + 2)
        for cur, next in dislikes:
            if uf.isConnected(cur, next):
                return False
            uf.union(cur, next + n)
            uf.union(cur + n, next)
        return True

