from collections import defaultdict
from typing import Callable, DefaultDict, List, Optional, Tuple


class WeightedUnionFind:
    """维护分量和的并查集."""

    __slots__ = "_parent", "_value", "_delta", "_total", "_part"

    def __init__(self, n: int):
        self._parent = [-1] * n
        self._value = [0] * n
        self._delta = [0] * n
        self._total = [0] * n
        self._part = n

    def add(self, u: int, delta: int) -> None:
        """u的值加上delta."""
        self._value[u] += delta
        self._total[self.find(u)] += delta

    def addGroup(self, u: int, delta: int) -> None:
        """u所在集合的值加上delta."""
        root = self.find(u)
        self._delta[root] += delta
        self._total[root] -= self._parent[root] * delta

    def get(self, u: int) -> int:
        """u的值."""
        return self._value[u] + self._find(u)[1]

    def getGroup(self, u: int) -> int:
        """u所在集合的值."""
        return self._total[self.find(u)]

    def union(self, u: int, v: int, f: Optional[Callable[[int, int], None]] = None) -> bool:
        u = self.find(u)
        v = self.find(v)
        if u == v:
            return False
        if self._parent[u] > self._parent[v]:
            u, v = v, u
        self._parent[u] += self._parent[v]
        self._parent[v] = u
        self._delta[v] -= self._delta[u]
        self._total[u] += self._total[v]
        self._part -= 1
        if f:
            f(u, v)
        return True

    def find(self, u: int) -> int:
        return self._find(u)[0]

    def isConnected(self, u: int, v: int) -> bool:
        return self.find(u) == self.find(v)

    def getSize(self, u: int) -> int:
        return -self._parent[self.find(u)]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        groups = defaultdict(list)
        for i in range(len(self._parent)):
            root = self.find(i)
            if root in groups:
                groups[root].append(i)
            else:
                groups[root] = [i]
        return groups

    @property
    def part(self) -> int:
        return self._part

    def _find(self, u: int) -> Tuple[int, int]:
        if self._parent[u] < 0:
            return u, self._delta[u]
        p = self._find(self._parent[u])
        first = p[0]
        second = p[1] + self._delta[u] - self._delta[first]
        self._parent[u] = first
        self._delta[u] = second
        return p


if __name__ == "__main__":
    N = 5
    uf = WeightedUnionFind(N)
    for i in range(N):
        uf.add(i, i)
    uf.union(2, 4)
    for i in range(N):
        print(uf.get(i), uf.getGroup(i))
    print()
    uf.union(3, 4)
    for i in range(N):
        print(uf.get(i), uf.getGroup(i))
    print()
    uf.add(3, 10)
    for i in range(N):
        print(uf.get(i), uf.getGroup(i))
    print()
    uf.addGroup(4, 5)
    for i in range(N):
        print(uf.get(i), uf.getGroup(i))
