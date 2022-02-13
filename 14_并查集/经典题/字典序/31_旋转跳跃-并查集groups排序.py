# 元素是0-n-1的并查集写法，不支持动态添加
from collections import defaultdict
from typing import Dict, Generic, Iterable, List, TypeVar


T = TypeVar('T')


# 当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加
class UnionFindMap(Generic[T]):
    def __init__(self, iterable: Iterable[T] = None):
        self.count = 0
        self.parent = dict()
        self.rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self._add(item)

    def union(self, key1: T, key2: T) -> bool:
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
            self._add(key)
            return key

        while self.parent.get(key, key) != key:
            self.parent[key] = self.parent[self.parent[key]]
            key = self.parent[key]
        return key

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroup(self) -> Dict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def _add(self, key: T) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.count += 1
        return True


class Point:
    def __init__(self, a=0, b=0):
        self.x = a
        self.y = b


class Solution:
    def solve(self, n, m, perm, Pair):
        """返回字典序最小的排列"""
        uf = UnionFindMap(list(range(1, n + 1)))
        for p in Pair:
            x, y = p.x, p.y
            uf.union(x, y)
        group = uf.getGroup()
        for points in group.values():
            points.sort(key=lambda x: -perm[x - 1])

        res = []
        for i in range(1, n + 1):
            root = uf.find(i)
            index = group[root].pop()
            res.append(perm[index - 1])
        return res

