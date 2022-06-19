from collections import defaultdict
from itertools import product
from typing import DefaultDict, Generic, Iterable, List, Optional, TypeVar


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


# 判断该场地是否可以保证每组摄像头至少有一个负责人管理。
class Solution:
    def isCompliance(self, distance: List[List[int]], k: int) -> bool:
        uf = UnionFindMap()
        for i, j in product(range(len(distance)), repeat=2):
            if distance[i][j] <= 2:
                uf.union(i, j)
        return len(uf.getRoots()) <= k


print(Solution().isCompliance([[0, 3, 3], [3, 0, 3], [3, 3, 0]], 2))
print(
    Solution().isCompliance(
        [
            [0, 4, 3, 2, 4, 4, 2, 3, 4, 4],
            [4, 0, 5, 1, 4, 4, 2, 4, 5, 6],
            [3, 5, 0, 3, 3, 3, 1, 6, 4, 3],
            [2, 1, 3, 0, 3, 6, 3, 3, 1, 4],
            [4, 4, 3, 3, 0, 1, 3, 1, 4, 4],
            [4, 4, 3, 6, 1, 0, 4, 7, 5, 6],
            [2, 2, 1, 3, 3, 4, 0, 5, 4, 4],
            [3, 4, 6, 3, 1, 7, 5, 0, 4, 3],
            [4, 5, 4, 1, 4, 5, 4, 4, 0, 5],
            [4, 6, 3, 4, 4, 6, 4, 3, 5, 0],
        ],
        3,
    )
)
# False True
3
