from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMap(Generic[T]):
    """当元素不是数组index时(例如字符串)，更加通用的并查集写法，支持动态添加"""

    __slots__ = ("part", "_parent", "_rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self._parent = dict()
        self._rank = dict()
        for item in iterable or []:
            self.add(item)

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self._rank[root1] > self._rank[root2]:
            root1, root2 = root2, root1
        self._parent[root1] = root2
        self._rank[root2] += self._rank[root1]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        if key not in self._parent:
            self.add(key)
            return key

        while self._parent.get(key, key) != key:
            self._parent[key] = self._parent[self._parent[key]]
            key = self._parent[key]
        return key

    def isConnected(self, key1: T, key2: T) -> bool:
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self._parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getSize(self, key: T) -> int:
        return self._rank[self.find(key)]

    def add(self, key: T) -> bool:
        if key in self._parent:
            return False
        self._parent[key] = key
        self._rank[key] = 1
        self.part += 1
        return True

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self._parent


class UnionFindArray:
    """元素是0-n-1的并查集写法,不支持动态添加

    初始化的连通分量个数 为 n
    """

    __slots__ = ("n", "part", "_parent", "_rank")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self._rank = [1] * n

    def find(self, x: int) -> int:
        while self._parent[x] != x:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
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
        return list(set(self.find(key) for key in self._parent))

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


class UnionFindMap2(Generic[T]):
    """不自动合并 需要手动add添加元素"""

    __slots__ = ("part", "_parent", "_rank")

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self._parent = dict()
        self._rank = defaultdict(lambda: 1)
        for item in iterable or []:
            self.add(item)

    def add(self, key: T) -> "UnionFindMap2[T]":
        if key in self._parent:
            return self
        self._parent[key] = key
        self._rank[key] = 1
        self.part += 1
        return self

    def union(self, key1: T, key2: T) -> bool:
        """rank一样时 默认key2作为key1的父节点"""
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2 or root1 not in self._parent or root2 not in self._parent:
            return False
        if self._rank[root1] > self._rank[root2]:
            root1, root2 = root2, root1
        self._parent[root1] = root2
        self._rank[root2] += self._rank[root1]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        """此处不自动add"""
        if key not in self._parent:
            return key

        if key != self._parent[key]:
            root = self.find(self._parent[key])
            self._parent[key] = root
        return self._parent[key]

    def isConnected(self, key1: T, key2: T) -> bool:
        if key1 not in self._parent or key2 not in self._parent:
            return False
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self._parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self._parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getSize(self, key: T) -> int:
        return self._rank[self.find(key)]

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self._parent


class UnionFindGraph:
    """并查集维护无向图每个连通块的边数和顶点数"""

    __slots__ = ("n", "part", "_parent", "vertex", "edge")

    def __init__(self, n: int):
        self.n = n
        self.part = n
        self._parent = list(range(n))
        self.vertex = [1] * n  # 每个联通块的顶点数
        self.edge = [0] * n  # 每个联通块的边数

    def find(self, x: int) -> int:
        while x != self._parent[x]:
            self._parent[x] = self._parent[self._parent[x]]
            x = self._parent[x]
        return x

    def union(self, x: int, y: int) -> bool:
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

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def getRoots(self) -> List[int]:
        return list(set(self.find(i) for i in range(self.n)))

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part


class ATCRevocableUnionFindArray:
    """维护分量之和的可撤销并查集"""

    __slots__ = ("n", "parentSize", "sum", "history")

    def __init__(self, n: int):
        self.n = n
        self.parentSize = [-1] * n
        self.sum = [0] * n
        self.history = []

    def addSum(self, i: int, delta: int):
        """第i个元素的值加上delta"""
        x = i
        while x >= 0:
            self.sum[x] += delta
            x = self.parentSize[x]

    def union(self, a: int, b: int) -> bool:
        x = self.find(a)
        y = self.find(b)
        if -self.parentSize[x] < -self.parentSize[y]:
            x, y = y, x
        self.history.append((x, self.parentSize[x]))
        self.history.append((y, self.parentSize[y]))
        if x == y:
            return False
        self.parentSize[x] += self.parentSize[y]
        self.parentSize[y] = x
        self.sum[x] += self.sum[y]
        return True

    def find(self, a: int) -> int:
        x = a
        while self.parentSize[x] >= 0:
            x = self.parentSize[x]
        return x

    def isConnected(self, a: int, b: int) -> bool:
        return self.find(a) == self.find(b)

    def revocate(self) -> bool:
        if not self.history:
            return False
        y, py = self.history.pop()
        x, px = self.history.pop()
        if self.parentSize[x] != px:
            self.sum[x] -= self.sum[y]
        self.parentSize[x] = px
        self.parentSize[y] = py
        return True

    def getComponentSum(self, i: int) -> int:
        return self.sum[self.find(i)]
