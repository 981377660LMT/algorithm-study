from collections import Counter
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# N 頂点
# M 辺の単純無向グラフが与えられます。頂点には
# 1,2,…,N の番号が、辺には
# 1,2,…,M の番号が付けられています。
# 辺
# i(i=1,2,…,M) は頂点
# u
# i
# ​
#  ,v
# i
# ​
#   を結んでいます。

# このグラフがパスグラフであるか判定してください。

# パスグラフとは
# 頂点に
# 1,2,…,N の番号が付けられた
# N 頂点のグラフがパスグラフであるとは、
# (1,2,…,N) を並べ変えて得られる数列
# (v
# 1
# ​
#  ,v
# 2
# ​
#  ,…,v
# N
# ​
#  ) であって、以下の条件を満たすものが存在することをいいます。
# 全ての
# i=1,2,…,N−1 に対して、頂点
# v
# i
# ​
#  ,v
# i+1
# ​
#   を結ぶ辺が存在する
# 整数
# i,j が
# 1≤i,j≤N,∣i−j∣≥2 を満たすならば、頂点
# v
# i
# ​
#  ,v
# j
# ​
#   を結ぶ辺は存在しない
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
    if m != n - 1:
        print("No")
        exit(0)
    uf = UnionFindArray(n)
    for _ in range(m):
        u, v = map(int, input().split())
        u -= 1
        v -= 1
        if uf.isConnected(u, v):
            print("No")
            exit(0)
        uf.union(u, v)
    print("Yes")
