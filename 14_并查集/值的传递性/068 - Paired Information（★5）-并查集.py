# q个条件
# Ti=0 时 有关系 xi+1=yi
# TI=1 时 假定xi=vi 判断yi是否有确定的值
# 题目不会给出矛盾的数据
# 如果不确定 输出 'Ambiguous'
# N,Q<=1e5
# Xi,Yi<=2e9

# 并查集维护距离???

import sys
from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = int(1e9 + 7)


T = TypeVar('T', bound=Hashable)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
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
        self.part -= 1
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

    def getGroups(self) -> DefaultDict[T, List[T]]:
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
        self.part += 1
        return True

    def __str__(self) -> str:
        return '\n'.join(f'{root}: {member}' for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


n = int(input())
q = int(input())
uf = UnionFindMap[int]()
queries = []
values = [-1] * (n + 1)
for _ in range(q):
    t, x, y, v = map(int, input().split())
    queries.append((t, x, y, v))
    if t == 0:
        values[y] = v
