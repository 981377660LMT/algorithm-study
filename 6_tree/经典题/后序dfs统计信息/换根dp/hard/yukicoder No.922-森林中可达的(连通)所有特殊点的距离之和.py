# https://yukicoder.me/problems/no/922

# yukicoder No.922-森林修建机场
# 给定一颗树,切断一些边，保留其中的m条边.
# 有q个移动,每次移动从ai到bi
# 为了完成移动,决定在每个连通分量设置一个飞机场
# 飞机场之间距离为0
# !合理安排飞机场的位置，使得移动的总距离最小，求出最小距离.

# 1. 如果移动a,b在同一个连通分量 =>  lca 求距离
# 2. 如果移动a,b不在同一个连通分量 => 移动到飞机场


from LCA import LCA_HLD
from RerootingForest import RerootingForest

import sys
from typing import Tuple, List, DefaultDict
from collections import defaultdict

sys.setrecursionlimit(int(1e7))
input = lambda: sys.stdin.readline().rstrip("\r\n")
INF = int(4e18)


class _UF:

    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
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


if __name__ == "__main__":

    n, m, q = map(int, input().split())
    # 树中保留的边
    edges = list((a - 1, b - 1) for a, b in tuple(map(int, input().split()) for _ in range(m)))
    queries = list((a - 1, b - 1) for a, b in tuple(map(int, input().split()) for _ in range(q)))

    tree = [[] for _ in range(n)]
    uf = _UF(n)
    LCA = LCA_HLD(n)
    for u, v in edges:
        tree[u].append((v, 1))
        tree[v].append((u, 1))
        uf.union(u, v)
        LCA.addEdge(u, v, 1)
    LCA.build(root=-1)

    res = 0
    weights = [0] * n  # 特殊点的权重
    for u, v in queries:
        if uf.isConnected(u, v):
            res += LCA.dist(u, v)
        else:
            weights[u] += 1
            weights[v] += 1

    """求`原树`中每个点到可达的(连通)所有特殊点的距离之和."""

    E = Tuple[int, int]  # (distSum,size) 树中距离之和, 子树中特殊点的大小

    def e(root: int) -> E:
        return (0, 0)

    def op(childRes1: E, childRes2: E) -> E:
        dist1, size1 = childRes1
        dist2, size2 = childRes2
        return (dist1 + dist2, size1 + size2)

    def composition(fromRes: E, parent: int, cur: int, direction: int) -> E:
        """direction: 0: cur -> parent, 1: parent -> cur"""
        dist, size = fromRes
        from_ = cur if direction == 0 else parent
        count = weights[from_]
        return (dist + size + count, size + count)

    R = RerootingForest(n)
    for u, v in edges:
        R.addEdge(u, v)
    groups = uf.getGroups()
    for root, nodes in groups.items():
        dp = R.rerooting(e, op, composition, groupRoot=root)
        values = dp.values()
        res += min(d for d, _ in values)
    print(res)
