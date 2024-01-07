# IntervalUnion/UnionInterval/IntervalGraphUnion

from collections import defaultdict
from typing import Callable, DefaultDict, List, Tuple


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


def unionIntervals(intervals: List[Tuple[int, int]]) -> "UnionFindArray":
    """
    给定n个区间[start,end), 返回合并后的并查集.
    """
    n = len(intervals)
    order = sorted(range(n), key=lambda i: (intervals[i][0], -intervals[i][1]))

    uf = UnionFindArray(n)
    keep = []
    for j in order:
        if keep:
            i = keep[-1]
            startI, endI = intervals[i]
            startJ, endJ = intervals[j]
            if endJ <= endI and endJ - startJ < endI - startI:
                uf.union(i, j)
                continue
        keep.append(j)

    for k in range(len(keep) - 1):
        i, j = keep[k], keep[k + 1]
        startI, endI = intervals[i]
        startJ, endJ = intervals[j]
        if max(startI, startJ) < min(endI, endJ):
            uf.union(i, j)

    return uf


if __name__ == "__main__":
    intervals = [(1, 3), (2, 4), (5, 6)]
    uf = unionIntervals(intervals)
    print(uf)
