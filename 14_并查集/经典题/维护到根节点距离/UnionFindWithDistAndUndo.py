from collections import defaultdict
from typing import Callable, DefaultDict, Generic, List, Tuple, TypeVar


D = TypeVar("D")


class UnionFindWithDistAndUndo(Generic[D]):
    """
    维护到每个组根节点距离的可撤销并查集.
    用于维护环的权值，树上的距离等.
    """

    __slots__ = ("_data", "_history", "_n", "_e", "_op", "_inv")

    def __init__(self, n: int, e: Callable[[], D], op: Callable[[D, D], D], inv: Callable[[D], D]):
        self._data: List[Tuple[int, D]] = [(-1, e()) for _ in range(n)]
        self._history = []
        self._n = n
        self._e = e
        self._op = op
        self._inv = inv

    def union(self, parent: int, child: int, dist: D) -> bool:
        """
        distToRoot(parent) = distToRoot(child) + dist.
        如果组内两点距离存在矛盾(沿着不同边走距离不同),返回false.
        """
        v1, x1 = self.find(parent)
        v2, x2 = self.find(child)
        if v1 == v2:
            return dist == self._e()
        s1, s2 = -self._data[v1][0], -self._data[v2][0]
        if s1 < s2:
            v1, v2 = v2, v1
            x1, x2 = x2, x1
            dist = self._inv(dist)
        # v1 <- v2
        dist = self._op(x1, dist)
        dist = self._op(dist, self._inv(x2))
        self._history.append((v2, self._data[v2]))
        self._data[v2] = (v1, dist)
        self._history.append((v1, self._data[v1]))
        self._data[v1] = (-(s1 + s2), self._e())
        return True

    def find(self, v: int) -> Tuple[int, D]:
        """返回v所在组的根节点和到v到根节点的距离."""
        root, distToRoot = v, self._e()
        while True:
            p, dist = self._data[root]
            if p < 0:
                break
            distToRoot = self._op(distToRoot, dist)
            root = p
        return root, distToRoot

    def dist(self, x: int, y: int) -> D:
        """返回x到y的距离`f(x) - f(y)`."""
        vx, dx = self.find(x)
        vy, dy = self.find(y)
        if vx != vy:
            raise ValueError("x and y are not in the same set")
        return self._op(dx, self._inv(dy))

    def distToRoot(self, x: int) -> D:
        return self.find(x)[1]

    def getTime(self) -> int:
        return len(self._history)

    def rollback(self, time: int) -> None:
        while len(self._history) > time:
            v, value = self._history.pop()
            self._data[v] = value

    def undo(self) -> bool:
        if not self._history:
            return False
        v, value = self._history.pop()
        self._data[v] = value
        return True

    def getSize(self, x: int) -> int:
        root, _ = self.find(x)
        return -self._data[root][0]

    def getGroups(self) -> DefaultDict[int, List[int]]:
        res = defaultdict(list)
        for i in range(self._n):
            root, _ = self.find(i)
            res[root].append(i)
        return res


if __name__ == "__main__":
    uf = UnionFindWithDistAndUndo(5, lambda: 0, lambda x, y: x + y, lambda x: -x)
    print(uf.union(0, 1, 1))
    print(uf.union(2, 3, 2))
    print(uf.union(0, 2, 3))
    print(uf.union(0, 4, 4))
    print(uf.find(4))
