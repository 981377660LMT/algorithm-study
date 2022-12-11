"""
给定一棵树,f(u,v)表示u到v的最短路中边权的最大值。
对所有组合(i,j),求f(i,j)的和。
n<=1e5

!树上的并查集 边权排序+离线查询
顺序遍历边,然后并查集合并 (每条边作为最大值的贡献)
`res += w * uf.rank[uf.find(u)] * uf.rank[uf.find(v)]`
"""

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

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


if __name__ == "__main__":
    n = int(input())
    edges = []
    for _ in range(n - 1):
        u, v, w = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v, w))

    edges.sort(key=lambda x: x[2])
    uf = UnionFindArray(n)
    res = 0
    for u, v, w in edges:
        res += w * uf.rank[uf.find(u)] * uf.rank[uf.find(v)]
        uf.union(u, v)
    print(res)
