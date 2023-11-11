# 有一个无向连通图 M，每次询问给出一条边 e，
# 询问 e 加入到 M 中是否会影响 M 的最小生成树。
# !把 M 中的边和查询边放到一起跑最小生成树算法，但是注意查询边只判断不改变连通性。


import sys
from collections import defaultdict
from typing import DefaultDict, List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        if x != self.parent[x]:
            self.parent[x] = self.find(self.parent[x])
        return self.parent[x]

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
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


n, m, q = map(int, input().split())
edges1 = []
for _ in range(m):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    edges1.append((u, v, w, -1))  # 最后一个i用来是识别边的种类

edges2 = []
for i in range(q):
    u, v, w = map(int, input().split())
    u, v = u - 1, v - 1
    edges2.append((u, v, w, i))

res = [False] * q
edges = sorted(edges1 + edges2, key=lambda x: x[2])
uf = UnionFindArray(n)
for u, v, w, id in edges:
    root1, root2 = uf.find(u), uf.find(v)
    if root1 != root2:
        if id == -1:
            uf.union(root1, root2)
        else:
            res[id] = True

for flag in res:
    print("Yes" if flag else "No")
