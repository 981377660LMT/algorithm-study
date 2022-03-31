# 当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加
from collections import defaultdict
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

        if key1 not in self.parent or key2 not in self.parent:
            return False

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


class Solution:
    def solve(self, sx, sy, ex, ey, roads):
        # 从起点不断加边，直到起点和终点联通或者 roads 加入完了
        start = (sx, sy)
        end = (ex, ey)
        # 注意特判
        for dx, dy in [(0, 0), (1, 0), (-1, 0), (0, 1), (0, -1)]:
            x = sx + dx
            y = sy + dy
            if (x, y) == (ex, ey):
                return 0

        # 注意不在并查集里的元素，不可以自动加入
        uf = UnionFindMap()
        uf.add(start)
        uf.add(end)
        for i, road in enumerate(map(tuple, roads)):
            uf.add(road)
            for dx, dy in [(1, 0), (-1, 0), (0, 1), (0, -1)]:
                x = road[0] + dx
                y = road[1] + dy
                uf.union(road, (x, y))
                if uf.isConnected(start, end):
                    return i + 1

        return -1
