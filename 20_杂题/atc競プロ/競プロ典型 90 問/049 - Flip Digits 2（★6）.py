# 商店里有一些商品 可以反转[left,right]这一段子串 (1=>0 0=>1)
# 起始时 有一个长为n的全0字符串
# !问如何购买商品 使得使用购买的商品可以得到任意01排列(2^n个) 且花费最少
# 求出最少花费 如果不行 则输出-1
# n,m<=1e5

# 1.条件变形:
# !反转[left,right]这一段子串等价于
# !反转前缀 pre[left-1]后再反转前缀pre[right]
# !因此可以在left-1 => right 连一条无向边
# !可以得到任意01排列 等价于n+1个顶点全部联通
# 2. 最小生成树

import sys
from collections import defaultdict
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


sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline

n, m = map(int, input().split())
edges = []
for _ in range(m):
    cost, left, right = map(int, input().split())  # 1 <= left <= right <= n
    left -= 1
    edges.append((left, right, cost))

res = kruskal(n + 1, edges)

print(-1 if res == int(1e20) else res)
