# TODO
# 看看哪个好/或者是hitonanode的

from collections import defaultdict
from typing import DefaultDict, List


class SelectOneFromEachPair:
    """
    维护是树的连通块(无环)的连通块个数.
    !可撤销并查集.
    """

    __slots__ = ("_optStack", "_parent", "_treeCount", "part", "vertex", "edge")

    def __init__(self) -> None:
        self._optStack = []
        self._parent = dict()
        self._treeCount = 0  # !所有连通分量min(顶点数,边数)的和
        self.part = 0
        self.vertex = dict()
        self.edge = dict()

    def add(self, key: T) -> bool:
        if key in self._parent:
            return False
        self._parent[key] = key
        self._treeCount += 1
        self.vertex[key] = 1
        self.edge[key] = 0
        self.part += 1
        return True

    def union(self, x: T, y: T) -> bool:
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            self.edge[rootX] += 1  # !两个顶点已经在同一个连通块了，这个连通块的边数+1
            return False
        if self.vertex[rootX] > self.vertex[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self.vertex[rootY] += self.vertex[rootX]
        self.edge[rootY] += self.edge[rootX] + 1
        self.part -= 1
        return True

    def find(self, x: T) -> T:
        if x not in self._parent:
            self.add(x)
            return x
        while self._parent.get(x, x) != x:
            x = self._parent[x]
        return x

    def countTree(self) -> int:
        return self._treeCount

    def revocate(self) -> None:
        """合并失败时也要撤销."""
        ...

    def getSize(self, key: T) -> int:
        return self.vertex[self.find(key)]

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def isConnected(self, x: T, y: T) -> bool:
        return self.find(x) == self.find(y)

    def __len__(self) -> int:
        return len(self._parent)
