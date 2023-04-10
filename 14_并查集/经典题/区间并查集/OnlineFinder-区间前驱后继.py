# 寻找前驱后继/区间删除

from typing import Optional


class OnlineFinder:
    """利用并查集寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
    初始时,所有位置都未被访问过.
    """

    __slots__ = "_n", "_lParent", "_rParent", "_removed"

    def __init__(self, n: int):
        self._n = n
        self._lParent = list(range(n + 2))  # 0 和 n + 1 为哨兵, 实际使用[1,n]
        self._rParent = list(range(n + 2))

    def prev(self, x: int) -> Optional[int]:
        """找到x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        res = self._lFind(x + 1)
        return res - 1 if res != 0 else None

    def next(self, x: int) -> Optional[int]:
        """x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        res = self._rFind(x + 1)
        return res - 1 if res != self._n + 1 else None

    def erase(self, left: int, right=-1) -> None:
        """删除[left, right)区间内的元素.
        0<=left<=right<=n.
        """
        if right == -1:
            right = left + 1
        if left >= right:
            return
        left, right = left + 1, right + 1

        leftRoot, rightRoot = self._rFind(left), self._rFind(right)
        while rightRoot != leftRoot:
            self._rUnion(leftRoot, leftRoot + 1)
            leftRoot = self._rFind(leftRoot + 1)

        leftRoot, rightRoot = self._lFind(left - 1), self._lFind(right - 1)
        while rightRoot != leftRoot:
            self._lUnion(rightRoot, rightRoot - 1)
            rightRoot = self._lFind(rightRoot - 1)

    def _lUnion(self, x: int, y: int) -> None:
        if x < y:
            x, y = y, x
        rootX = self._lFind(x)
        rootY = self._lFind(y)
        if rootX == rootY:
            return
        self._lParent[rootX] = rootY

    def _rUnion(self, x: int, y: int) -> None:
        if x > y:
            x, y = y, x
        rootX = self._rFind(x)
        rootY = self._rFind(y)
        if rootX == rootY:
            return
        self._rParent[rootX] = rootY

    def _lFind(self, x: int) -> int:
        while x != self._lParent[x]:
            self._lParent[x] = self._lParent[self._lParent[x]]
            x = self._lParent[x]
        return x

    def _rFind(self, x: int) -> int:
        while x != self._rParent[x]:
            self._rParent[x] = self._rParent[self._rParent[x]]
            x = self._rParent[x]
        return x


if __name__ == "__main__":

    def pre(pos: int):
        return next((i for i in range(pos, -1, -1) if ok[i]), None)

    def nxt(pos: int):
        return next((i for i in range(pos, n) if ok[i]), None)

    def erase(left: int, right: int):
        for i in range(left, right):
            ok[i] = False

    from random import randint

    for _ in range(100):
        n = randint(1, 100)
        F = OnlineFinder(n)
        ok = [True] * n
        for _ in range(100):
            left, right = randint(0, n), randint(0, n)
            F.erase(left, right)
            erase(left, right)
            for i in range(n):
                assert F.prev(i) == pre(i), (i, F.prev(i), pre(i))
                assert F.next(i) == nxt(i), (i, F.next(i), nxt(i))
    print("Done!")
