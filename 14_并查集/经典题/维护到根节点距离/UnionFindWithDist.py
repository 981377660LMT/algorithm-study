"""
带权并查集(维护到每个组根节点距离的并查集)

- 注意距离是`有向`的
  例如维护和距离的并查集时,a->b 的距离是正数,b->a 的距离是负数
- 如果组内两点距离存在矛盾(沿着不同边走距离不同),那么在组内会出现正环
"""

# Class:
#  UnionFindArrayWithDist1(n) : 数组实现的并查集,距离为加法.
#  UnionFindArrayWithDist2(n) : 数组实现的并查集,距离为乘法.
#  UnionFindMapWithDist1() : 字典实现的并查集,距离为加法.
#  UnionFindMapWithDist2() : 字典实现的并查集,距离为乘法.


# API:
#  Union(x,y,dist) : p(x) = p(y) + dist. 如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
#  Find(x) : 返回x所在组的根节点.
#  Dist(x,y) : 返回x到y的距离.
#  DistToRoot(x) : 返回x到所在组根节点的距离.

from collections import defaultdict
from typing import Callable, DefaultDict, List, Optional


class UnionFindArrayWithDist1:
    """维护到根节点距离的并查集.距离为加法."""

    __slots__ = ("part", "_data", "_potential")

    def __init__(self, n: int):
        self.part = n
        self._data = [-1] * n
        self._potential = [0] * n

    def union(
        self, x: int, y: int, dist: int, cb: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """
        p(x) = p(y) + dist.
        如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
        """
        dist += self.distToRoot(y) - self.distToRoot(x)
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
            cb(y, x)
        return True

    def find(self, x: int) -> int:
        if self._data[x] < 0:
            return x
        r = self.find(self._data[x])
        self._potential[x] += self._potential[self._data[x]]
        self._data[x] = r
        return r

    def dist(self, x: int, y: int) -> int:
        """返回x到y的距离`f(x) - f(y)`."""
        return self.distToRoot(x) - self.distToRoot(y)

    def distToRoot(self, x: int) -> int:
        """返回x到所在组根节点的距离`f(x) - f(find(x))`."""
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


class UnionFindArrayWithDist2:
    """维护到根节点距离的并查集.距离为乘法."""

    __slots__ = ("part", "_data", "_potential")

    def __init__(self, n: int):
        self.part = n
        self._data = [-1] * n
        self._potential = [1.0] * n

    def union(
        self, x: int, y: int, dist: float, cb: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """
        p(x) = p(y) * dist.
        如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
        """
        dist *= self.distToRoot(y) / self.distToRoot(x)
        x, y = self.find(x), self.find(y)
        if x == y:
            return dist == 1
        if self._data[x] < self._data[y]:
            x, y = y, x
            dist = 1 / dist
        self._data[y] += self._data[x]
        self._data[x] = y
        self._potential[x] = dist
        self.part -= 1
        if cb is not None:
            cb(y, x)
        return True

    def find(self, x: int) -> int:
        if self._data[x] < 0:
            return x
        r = self.find(self._data[x])
        self._potential[x] *= self._potential[self._data[x]]
        self._data[x] = r
        return r

    def dist(self, x: int, y: int) -> float:
        """返回x到y的距离`f(x)/f(y)`."""
        return self.distToRoot(x) / self.distToRoot(y)

    def distToRoot(self, x: int) -> float:
        """返回x到所在组根节点的距离`f(x)/f(find(x))`."""
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


def id(o: object) -> int:
    if o not in _pool:
        _pool[o] = len(_pool)
    return _pool[o]


_pool = dict()


class UnionFindMapWithDist1:
    """维护到根节点距离的并查集.距离为加法."""

    __slots__ = ("part", "_data", "_potential")

    def __init__(self):
        self.part = 0
        self._data = dict()
        self._potential = dict()

    def union(
        self, x: int, y: int, dist: int, cb: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """
        p(x) = p(y) + dist.
        !如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
        """
        dist += self.distToRoot(y) - self.distToRoot(x)
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
            cb(y, x)
        return True

    def find(self, x: int) -> int:
        if x not in self._data:
            self.add(x)
            return x
        if self._data[x] < 0:
            return x
        r = self.find(self._data[x])
        self._potential[x] += self._potential[self._data[x]]
        self._data[x] = r
        return r

    def dist(self, x: int, y: int) -> int:
        """返回x到y的距离`f(x) - f(y)`."""
        return self.distToRoot(x) - self.distToRoot(y)

    def distToRoot(self, x: int) -> int:
        """返回x到所在组根节点的距离`f(x) - f(find(x))`."""
        self.find(x)
        return self._potential[x]

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        res = defaultdict(list)
        for k in self._data:
            res[self.find(k)].append(k)
        return res

    def add(self, x: int) -> "UnionFindMapWithDist1":
        if x not in self._data:
            self._data[x] = -1
            self._potential[x] = 0
            self.part += 1
        return self

    def __contains__(self, x: int) -> bool:
        return x in self._data

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


class UnionFindMapWithDist2:
    """维护到根节点距离的并查集.距离为乘法."""

    __slots__ = ("part", "_data", "_potential")

    def __init__(self):
        self.part = 0
        self._data = dict()
        self._potential = dict()

    def union(
        self, x: int, y: int, dist: float, cb: Optional[Callable[[int, int], None]] = None
    ) -> bool:
        """
        p(x) = p(y) * dist.
        !如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
        """
        dist *= self.distToRoot(y) / self.distToRoot(x)
        x, y = self.find(x), self.find(y)
        if x == y:
            return dist == 1
        if self._data[x] < self._data[y]:
            x, y = y, x
            dist = 1 / dist
        self._data[y] += self._data[x]
        self._data[x] = y
        self._potential[x] = dist
        self.part -= 1
        if cb is not None:
            cb(y, x)
        return True

    def find(self, x: int) -> int:
        if x not in self._data:
            self.add(x)
            return x
        if self._data[x] < 0:
            return x
        r = self.find(self._data[x])
        self._potential[x] *= self._potential[self._data[x]]
        self._data[x] = r
        return r

    def dist(self, x: int, y: int) -> float:
        """返回x到y的距离`f(x)/f(y)`."""
        return self.distToRoot(x) / self.distToRoot(y)

    def distToRoot(self, x: int) -> float:
        """返回x到所在组根节点的距离`f(x)/f(find(x))`."""
        self.find(x)
        return self._potential[x]

    def isConnected(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return -self._data[self.find(x)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        res = defaultdict(list)
        for k in self._data:
            res[self.find(k)].append(k)
        return res

    def add(self, x: int) -> "UnionFindMapWithDist2":
        if x not in self._data:
            self._data[x] = -1
            self._potential[x] = 1.0
            self.part += 1
        return self

    def __contains__(self, x: int) -> bool:
        return x in self._data

    def __repr__(self) -> str:
        return "\n".join(f"{root}: {member}" for root, member in self.getGroups().items())


# https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=DSL_1_B&lang=ja
if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e6))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, q = map(int, input().split())
    uf = UnionFindMapWithDist1()
    for _ in range(q):
        op, *rest = map(int, input().split())
        if op == 0:
            x, y, w = rest
            uf.union(y, x, w)
        else:
            x, y = rest
            if not uf.isConnected(x, y):
                print("?")
            else:
                print(uf.dist(y, x))
