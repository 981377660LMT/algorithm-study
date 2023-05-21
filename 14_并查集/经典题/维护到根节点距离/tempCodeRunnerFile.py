"""
带权并查集(维护到每个组根节点距离的并查集)

- 注意距离是`有向`的
  例如维护和距离的并查集时,a->b 的距离是正数,b->a 的距离是负数
- 如果组内两点距离存在矛盾(沿着不同边走距离不同),那么在组内会出现正环
"""

from collections import defaultdict
from typing import Callable, DefaultDict, Generic, Hashable, Iterable, List, Optional, TypeVar


class UnionFindArrayWithDist:
    __slots__ = ("part", "_data", "_potential")

    def __init__(self, n: int):
        self.part = n
        self._data = [-1] * n
        self._potential = [0] * n

    def union(
        self, x: int, y: int, dist: int, cb: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        dist += self._potential[y] - self._potential[x]
        x, y = self.find(x), self.find(y)
        if x == y:
            return dist == 0
        if self._data[x] < self._data[y]:
            x, y = y, x
            dist = -dist
        self._data[y] += self._data[x]
        self._data[x] = y
        self._potential[x] = dist
        self.part -= 1
        if cb is not None:
            cb(x, y)
        return True

    def find(self, x: int) -> int:
        if self._data[x] < 0:
            return x
        r = self.find(self._data[x])
        self._potential[x] += self._potential[self._data[x]]
        self._data[x] = r
        return r

    def dist(self, x: int, y: int) -> int:
        return self.distToRoot(x) - self.distToRoot(y)

    def distToRoot(self, x: int) -> int:
        self.find(x)
        return self._potential[x]

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        res = defaultdict(list)
        for i in range(len(self._data)):
            res[self.find(i)].append(i)
        return res


T = TypeVar("T", bound=Hashable)


class UnionFindMapWithDist1(Generic[T]):
    """需要手动添加元素,维护乘积(距离)的并查集"""

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.distToRoot = defaultdict(lambda: 1.0)
        for item in iterable or []:
            self.add(item)

    def getDist(self, key1: T, key2: T) -> float:
        """有向边 key1 -> key2 的距离"""
        if (key1 not in self.parent) or (key2 not in self.parent):
            raise KeyError("key not in UnionFindMapWithDist")
        return self.distToRoot[key1] / self.distToRoot[key2]

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

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent


class UnionFindMapWithDist2(Generic[T]):
    """需要手动添加元素,维护加法(距离)的并查集"""

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self.part = 0
        self.parent = dict()
        self.distToRoot = defaultdict(int)
        for item in iterable or []:
            self.add(item)

    def dist(self, key1: T, key2: T) -> int:
        """有向边 key1 -> key2 的距离"""
        if (key1 not in self.parent) or (key2 not in self.parent):
            raise KeyError("key not in UnionFindMapWithDist")
        return self.distToRoot[key1] - self.distToRoot[key2]

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

    def getGroups(self) -> DefaultDict[T, List[T]]:
        groups = defaultdict(list)
        for key in self.parent:
            root = self.find(key)
            groups[root].append(key)
        return groups

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())

    def __len__(self) -> int:
        return self.part

    def __contains__(self, key: T) -> bool:
        return key in self.parent


# https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_1_B&lang=ja
if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, q = map(int, input().split())
    uf = UnionFindArrayWithDist(n)
    for _ in range(q):
        op, *rest = map(int, input().split())
        if op == 0:
            x, y, w = rest
            uf.union(x, y, w)
        else:
            x, y = rest
            if not uf.isConnected(x, y):
                print("?")
            else:
                print(uf.dist(x, y))
