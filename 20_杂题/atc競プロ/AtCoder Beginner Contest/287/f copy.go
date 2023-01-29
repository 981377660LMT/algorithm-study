import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# N 頂点の木があります。頂点には
# 1 から
# N までの番号が付いており、
# i 番目の辺は頂点
# a
# i
# ​
#   と頂点
# b
# i
# ​
#   を結んでいます。

# x=1,2,…,N に対して次の問題を解いてください。

# 木の頂点の部分集合
# V であって空でないものは
# 2
# N
#  −1 通り存在するが、そのうち
# V による誘導部分グラフの連結成分数が
# x であるようなものは何通りあるかを
# 998244353 で割った余りを求めよ。

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


# dp[i][v] 表示以 i 为根的子树中，连通块数为 v 的方案数

if __name__ == "__main__":
    n = int(input())
    adjList = [[] for _ in range(n)]
    uf = UnionFindArray(n)
    for _ in range(n - 1):
        a, b = map(int, input().split())
        adjList[a - 1].append(b - 1)
        adjList[b - 1].append(a - 1)
        uf.union(a - 1, b - 1)
