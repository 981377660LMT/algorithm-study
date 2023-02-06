import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


# N 頂点
# M 辺の単純無向グラフが与えられます。頂点には
# 1 から
# N の番号がついており、
# i 番目の辺は頂点
# A
# i
# ​
#   と頂点
# B
# i
# ​
#   を結んでいます。 このグラフから
# 0 本以上のいくつかの辺を削除してグラフが閉路を持たないようにするとき、削除する辺の本数の最小値を求めてください。

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
    n, m = map(int, input().split())
    uf = UnionFindArray(n)
    res = 0
    for _ in range(m):
        a, b = map(int, input().split())
        a, b = a - 1, b - 1
        if uf.isConnected(a, b):
            continue
        uf.union(a, b)
        res += 1
    print(m - res)
