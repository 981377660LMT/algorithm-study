from collections import defaultdict
from itertools import combinations
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, Tuple, TypeVar

T = TypeVar('T')


class UnionFindMap(Generic[T]):
    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.count = 0
        self.parent = dict()
        self.rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self.add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self.rank[root1] > self.rank[root2]:
            root1, root2 = root2, root1
        self.parent[root1] = root2
        self.rank[root2] += self.rank[root1]
        self.count -= 1
        return True

    def find(self, key: T) -> T:
        if key not in self.parent:
            self.add(key)
            return key

        while self.parent.get(key, key) != key:
            self.parent[key] = self.parent[self.parent[key]]
            key = self.parent[key]
        return key

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroup(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def add(self, key: T) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.count += 1
        return True


P = TypeVar('P', bound=Hashable)


def kruskal(n: int, edges: List[Tuple[P, P, int]]) -> int:
    """求最小生成树权值
    
    1. 边权排序
    2. 两两连接不连通的点
    """
    uf = UnionFindMap[P]()
    res, hit = 0, 0

    edges = sorted(edges, key=lambda e: e[2])
    for u, v, w in edges:
        root1, root2 = uf.find(u), uf.find(v)
        if root1 != root2:
            res += w
            uf.union(root1, root2)
            hit += 1

    return res if hit == n - 1 else int(1e20)


# 1 <= points.length <= 1000
class Solution:
    def minCostConnectPoints(self, points: List[List[int]]) -> int:
        n = len(points)
        edges = []
        P = [tuple(p) for p in points]
        for p1, p2 in combinations(P, 2):
            edges.append((p1, p2, abs(p1[0] - p2[0]) + abs(p1[1] - p2[1])))
        return kruskal(n, edges)

