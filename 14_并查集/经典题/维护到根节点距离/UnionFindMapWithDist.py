from collections import defaultdict
from typing import DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


T = TypeVar("T", bound=Hashable)


class UnionFindMapWithDist1(Generic[T]):
    """需要手动添加元素 维护乘积(距离)的并查集"""

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.distToRoot = defaultdict(lambda: 1.0)
        for item in iterable or []:
            self.add(item)

    def add(self, key: T) -> "UnionFindMapWithDist1[T]":
        if key in self.parent:
            return self
        self.parent[key] = key
        self.part += 1
        return self

    def union(self, son: T, father: T, dist: float) -> bool:
        """
        father 与 son 间的距离为 dist
        围绕着'到根的距离'进行计算
        注意从走两条路到新根节点的距离是一样的
        """
        root1 = self.find(son)
        root2 = self.find(father)
        if (root1 == root2) or (root1 not in self.parent) or (root2 not in self.parent):
            return False

        self.parent[root1] = root2
        # !1. 合并距离 保持两条路到新根节点的距离是一样的
        self.distToRoot[root1] = dist * self.distToRoot[father] / self.distToRoot[son]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        """此处不自动add"""
        if key not in self.parent:
            return key

        # !2. 从上往下懒更新到根的距离
        if key != self.parent[key]:
            root = self.find(self.parent[key])
            self.distToRoot[key] *= self.distToRoot[self.parent[key]]
            self.parent[key] = root
        return self.parent[key]

    def isConnected(self, key1: T, key2: T) -> bool:
        if (key1 not in self.parent) or (key2 not in self.parent):
            return False
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __str__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent


class UnionFindMapWithDist2(Generic[T]):
    """需要手动添加元素 维护加法(距离)的并查集"""

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.distToRoot = defaultdict(int)
        for item in iterable or []:
            self.add(item)

    def add(self, key: T) -> "UnionFindMapWithDist2[T]":
        if key in self.parent:
            return self
        self.parent[key] = key
        self.part += 1
        return self

    def union(self, son: T, father: T, dist: int) -> bool:
        """
        father 与 son 间的距离为 dist
        围绕着'到根的距离'进行计算
        注意从走两条路到新根节点的距离是一样的
        """
        root1 = self.find(son)
        root2 = self.find(father)
        if (root1 == root2) or (root1 not in self.parent) or (root2 not in self.parent):
            return False

        self.parent[root1] = root2
        # !1. 合并距离 保持两条路到新根节点的距离是一样的
        self.distToRoot[root1] = dist + self.distToRoot[father] - self.distToRoot[son]
        self.part -= 1
        return True

    def find(self, key: T) -> T:
        """此处不自动add"""
        if key not in self.parent:
            return key

        # !2. 从上往下懒更新到根的距离
        if key != self.parent[key]:
            root = self.find(self.parent[key])
            self.distToRoot[key] += self.distToRoot[self.parent[key]]
            self.parent[key] = root
        return self.parent[key]

    def isConnected(self, key1: T, key2: T) -> bool:
        if (key1 not in self.parent) or (key2 not in self.parent):
            return False
        return self.find(key1) == self.find(key2)

    def getRoots(self) -> List[T]:
        return list(set(self.find(key) for key in self.parent))

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __str__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent
