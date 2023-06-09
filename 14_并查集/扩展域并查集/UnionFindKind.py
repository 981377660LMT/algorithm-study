# 种类并查集 0:同类,1:不同类


from collections import defaultdict
from typing import DefaultDict, Iterable, List, Optional


class UnionFindKindArray:
    """数组实现的种类并查集."""

    __slots__ = ("_n", "_k", "_uf")

    def __init__(self, n: int, k=2) -> None:
        """
        n: 0-n-1.
        k: 元素之间的关系个数(关系种类数).

        例如:
        - 0表示同类, 1表示不同类.
        - 0表示x和y是同类,1表示x吃y,2表示x被y吃.
        """
        self._n = n
        self._k = k
        self._uf = _UArray(n * k)

    def union(self, a: int, b: int, relation: int) -> bool:
        """
        将a和b合并, 关系为relation.
        如果存在矛盾,返回False.否则返回True.
        """
        n, k = self._n, self._k
        for i in range(k):
            if i == relation:
                continue
            if self._uf.isSame(a, b + n * i):
                return False
        for i in range(k):
            self._uf.union(a + n * i, b + n * ((i + relation) % k))
        return True

    def hasRelation(self, a: int, b: int, relation: int) -> bool:
        """检查a和b是否满足关系relation."""
        n, k = self._n, self._k
        for i in range(k):
            if not self._uf.isSame(a + n * i, b + n * ((i + relation) % k)):
                return False
        return True

    def getGroups(self) -> DefaultDict[int, List[int]]:
        """返回每个种类的所有元素."""
        res = defaultdict(list)
        n = self._n
        for i in range(n):
            res[self._uf.find(i)].append(i)
        return res


class _UArray:
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

    def isSame(self, x: int, y: int) -> bool:
        return self.find(x) == self.find(y)

    def getSize(self, x: int) -> int:
        return self._rank[self.find(x)]


class UnionFindKindMap:
    """字典实现的种类并查集."""

    __slots__ = "_uf", "_n", "_k"

    def __init__(self, n: int, k=2) -> None:
        """
        n: 元素0-n-1.
        k: 元素之间的关系个数(关系种类数).

        例如:
        - 0表示同类, 1表示不同类.
        - 0表示x和y是同类,1表示x吃y,2表示x被y吃.
        """
        self._n = n
        self._k = k
        self._uf = _UMap()

    def union(self, a: int, b: int, relation: int) -> bool:
        """
        将a和b合并, 关系为relation.
        如果存在矛盾,返回False.否则返回True.
        """
        n, k = self._n, self._k
        for i in range(k):
            if i == relation:
                continue
            if self._uf.isSame(a, b + n * i):
                return False
        for i in range(k):
            self._uf.union(a + n * i, b + n * ((i + relation) % k))
        return True

    def hasRelation(self, a: int, b: int, relation: int) -> bool:
        """检查a和b是否满足关系relation."""
        n, k = self._n, self._k
        for i in range(k):
            if not self._uf.isSame(a + n * i, b + n * ((i + relation) % k)):
                return False
        return True

    def getGroups(self) -> DefaultDict[int, List[int]]:
        """返回每个种类的所有元素."""
        res = defaultdict(list)
        for key in self._uf.parent:
            if key < self._n:
                res[self._uf.find(key)].append(key)
        return res


class _UMap:
    __slots__ = ("part", "parent", "rank")

    def __init__(self, iterable: Optional[Iterable[int]] = None):
        self.part = 0
        self.parent = dict()
        self.rank = dict()
        for item in iterable or []:
            self.add(item)

    def union(self, key1: int, key2: int) -> bool:
        root1 = self.find(key1)
        root2 = self.find(key2)
        if root1 == root2:
            return False
        if self.rank[root1] > self.rank[root2]:
            root1, root2 = root2, root1
        self.parent[root1] = root2
        self.rank[root2] += self.rank[root1]
        self.part -= 1
        return True

    def find(self, key: int) -> int:
        if key not in self.parent:
            self.add(key)
            return key

        while self.parent.get(key, key) != key:
            self.parent[key] = self.parent[self.parent[key]]
            key = self.parent[key]
        return key

    def add(self, key: int) -> bool:
        if key in self.parent:
            return False
        self.parent[key] = key
        self.rank[key] = 1
        self.part += 1
        return True

    def isSame(self, key1: int, key2: int) -> bool:
        return self.find(key1) == self.find(key2)

    def getSize(self, x: int) -> int:
        return self.rank[self.find(x)]


if __name__ == "__main__":
    n = 4  # 1-4
    dislikes = [[1, 2], [1, 3], [2, 4]]
    uf = UnionFindKindArray(4, 2)
    for a, b in dislikes:
        uf.union(a - 1, b - 1, 1)
    print(uf.getGroups())
