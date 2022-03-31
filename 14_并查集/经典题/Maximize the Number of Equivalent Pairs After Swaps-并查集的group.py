# 当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加
from collections import defaultdict
from typing import Counter, DefaultDict, Generic, Iterable, List, Optional, TypeVar


T = TypeVar('T')


class UnionFindMap(Generic[T]):
    def __init__(self, iterable: Optional[Iterable[T]] = None):
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

    def getGroup(self) -> DefaultDict[T, List[T]]:
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


class Solution:
    def solve(self, A, B, C):
        """求交换后的最大相等对数"""
        n = len(A)
        uf = UnionFindMap(range(n))  # 存下标
        for u, v in C:
            uf.union(u, v)

        res = 0
        group = uf.getGroup()
        for root in group:
            # 可以互相交换的组
            indexes = group[root]
            counter = Counter(A[i] for i in indexes)
            # A的组里 看对应B的位置哪些可以放
            for i in indexes:
                if counter[B[i]]:
                    counter[B[i]] -= 1
                    res += 1
        return res


print(Solution().solve(A=[1, 2, 3, 4], B=[2, 1, 4, 3], C=[[0, 1], [2, 3]]))
