# E - Just one-无向图形成基环树的方案数
# !给定一个无向图，求有多少中安排方式使边有向并且每个点的出度为1(内向基环树森林)

# !并查集
# !对于每个连通块只需要判断当前是否存在且仅存在一个环，即点数和边数相等

# n,m<=2e5

import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
from collections import defaultdict
from typing import DefaultDict, List


class UnionFindGraph:
    """并查集维护每个连通块的边数和顶点数"""

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n  # 每个联通块的顶点数
        self.edge = [0] * n  # 每个联通块的边数

    def find(self, x: int) -> int:
        while x != self.parent[x]:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.edge[rootX] += 1  # !两个顶点已经在同一个连通块了，那么这个连通块的边数+1
            return False
        if self.rank[rootX] > self.rank[rootY]:
            rootX, rootY = rootY, rootX
        self.parent[rootX] = rootY
        self.rank[rootY] += self.rank[rootX]
        self.edge[rootY] += self.edge[rootX] + 1
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
        return list(set(self.find(i) for i in range(self.n)))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


if __name__ == "__main__":
    n, m = map(int, input().split())
    uf = UnionFindGraph(n)
    for _ in range(m):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        uf.union(u, v)

    res = 1
    for root in uf.getRoots():
        edge, vertex = uf.edge[root], uf.rank[root]
        if edge == vertex:
            res = res * 2 % MOD
        else:
            print(0)
            exit(0)
    print(res)
