from collections import defaultdict
from typing import Callable, DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


class UnionFindArraySimple:
    __slots__ = ("part", "n", "_data")

    def __init__(self, n: int):
        self.part = n
        self.n = n
        self._data = [-1] * n

    def union(
        self, key1: int, key2: int, beforeUnion: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        root1, root2 = self.find(key1), self.find(key2)
        if root1 == root2:
            return False
        if self._data[root1] > self._data[root2]:
            root1, root2 = root2, root1
        if beforeUnion is not None:
            beforeUnion(root1, root2)
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def unionTo(self, parent: int, child: int) -> bool:
        """定向合并, 将child合并到parent所在的连通分量中."""
        root1, root2 = self.find(parent), self.find(child)
        if root1 == root2:
            return False
        self._data[root1] += self._data[root2]
        self._data[root2] = root1
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        root = key
        while self._data[root] >= 0:
            root = self._data[root]
        while key != root:
            parent = self._data[key]
            self._data[key] = root
            key = parent
        return root

    def getSize(self, key: int) -> int:
        return -self._data[self.find(key)]


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
        """按秩合并."""
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

    def unionTo(self, child: T, parent: T) -> bool:
        """定向合并."""
        root1 = self.find(child)
        root2 = self.find(parent)
        if root1 == root2:
            return False
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
        """按秩合并."""
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

    def unionTo(self, child: int, parent: int) -> bool:
        """定向合并.将child的父节点设置为parent."""
        rootX = self.find(child)
        rootY = self.find(parent)
        if rootX == rootY:
            return False
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        return True

    def unionWithCallback(self, x: int, y: int, f: Callable[[int, int], None]) -> bool:
        """
        f: 合并后的回调函数, 入参为 (big, small)
        """
        rootX = self.find(x)
        rootY = self.find(y)
        if rootX == rootY:
            return False
        if self._rank[rootX] > self._rank[rootY]:
            rootX, rootY = rootY, rootX
        self._parent[rootX] = rootY
        self._rank[rootY] += self._rank[rootX]
        self.part -= 1
        f(rootY, rootX)
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
        self.vertex = [1] * n  # 每个联通块的顶点数
        self.edge = [0] * n  # 每个联通块的边数
        self._parent = list(range(n))

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

    def getSize(self, x: int) -> int:
        return self.vertex[self.find(x)]

    def getEdge(self, x: int) -> int:
        return self.edge[self.find(x)]

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for key in range(self.n):
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())
