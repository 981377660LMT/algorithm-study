# 无向图中x到y的路径是否唯一
# !等价于路径上所有边都是割边
# !并查集合并所有的割边顶点


from collections import defaultdict
from typing import DefaultDict, List
import sys

from Tarjan import Tarjan

sys.setrecursionlimit(int(1e9))
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


n = int(input())
adjList = [[] for _ in range(n)]

for _ in range(n):
    a, b = map(int, input().split())
    a, b = a - 1, b - 1
    adjList[a].append(b)
    adjList[b].append(a)

_, cuttingEdges = Tarjan.getCuttingPointAndCuttingEdge(n, adjList)
uf = UnionFindArray(n)
for a, b in cuttingEdges:
    uf.union(a, b)
q = int(input())
for _ in range(q):
    x, y = map(int, input().split())
    x, y = x - 1, y - 1
    if uf.isConnected(x, y):
        print("Yes")
    else:
        print("No")
