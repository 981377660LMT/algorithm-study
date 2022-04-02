from collections import defaultdict
from typing import DefaultDict, Generic, Iterable, List, Optional, TypeVar


# 当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加
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


# 不引爆炸弹的路径
class Solution:
    def solve(self, mines, width):
        isIntersected = (
            lambda i, j, d: (mines[i][0] - mines[j][0]) ** 2 + (mines[i][1] - mines[j][1]) ** 2
            <= d * d
        )

        n = len(mines)
        uf = UnionFindMap()
        left, right = -1, n
        for i in range(n):
            if mines[i][0] - mines[i][2] <= 0:
                uf.union(left, i)
            if width <= mines[i][0] + mines[i][2]:
                uf.union(right, i)
            for j in range(i + 1, n):
                if isIntersected(i, j, mines[i][2] + mines[j][2]):
                    uf.union(i, j)

        return not uf.isConnected(left, right)

