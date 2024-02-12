from collections import defaultdict
from typing import DefaultDict, List, Tuple


def kruskal(n: int, edges: List[Tuple[int, int, int]]) -> Tuple[int, List[bool], bool]:
    """
    Kruskal算法求无向图最小生成树(森林).

    返回值:
    - mstCost: 最小生成树(森林)的权值之和
    - inMst: 是否在最小生成树(森林)中
    - isTree: 是否是树
    """
    uf = UnionFindArray(n)
    count = 0
    mstCost, inMst, isTree = 0, [False] * len(edges), False
    order = sorted(range(len(edges)), key=lambda x: edges[x][2])
    for ei in order:
        u, v, w = edges[ei]
        if uf.union(u, v):
            inMst[ei] = True
            mstCost += w
            count += 1
            if count == n - 1:
                isTree = True
                break
    return mstCost, inMst, isTree


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    __slots__ = ("n", "part", "parent", "rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self.parent = list(range(n))
        self.rank = [1] * n

    def find(self, x: int) -> int:
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]
            x = self.parent[x]
        return x

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


if __name__ == "__main__":
    # https://www.luogu.com.cn/problem/P3366
    # P3366 【模板】最小生成树
    import sys

    input = lambda: sys.stdin.readline().rstrip("\r\n")
    n, m = map(int, input().split())
    edges = []
    for _ in range(m):
        u, v, w = map(int, input().split())
        edges.append((u - 1, v - 1, w))
    mstCost, _, isTree = kruskal(n, edges)
    print(mstCost if isTree else "orz")
