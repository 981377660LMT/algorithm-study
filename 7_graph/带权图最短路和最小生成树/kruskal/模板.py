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


# 1. 边权排序
# 2. 两两连接不连通的点
def kruskal(edges: List[List[int]]) -> int:
    """求最小生成树权值"""
    n = len(edges)
    uf = UnionFindArray(n + 10)
    res, hit = 0, 0

    # u,v,weight
    edges = sorted(edges, key=lambda e: e[2])
    for u, v, w in edges:
        root1, root2 = uf.find(u), uf.find(v)
        if root1 != root2:
            res += w
            uf.union(root1, root2)
            hit += 1

    return res if hit == n - 1 else int(1e20)
