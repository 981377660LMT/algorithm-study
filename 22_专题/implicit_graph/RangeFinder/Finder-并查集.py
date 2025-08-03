# 寻找前驱后继/区间删除

from typing import List, Optional


class Finder:
    """利用并查集寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
    初始时,所有位置都未被访问过.
    """

    ___slots___ = ("left", "right", "_n", "_data")

    def __init__(self, n: int):
        self._n = n
        n += 2
        self._data = [-1] * n  # 0 和 n + 1 为哨兵, 实际使用[1,n]
        self.left = list(range(n))  # 每个组的左边界
        self.right = [i + 1 for i in range(n)]  # 每个组的右边界

    def prev(self, x: int) -> Optional[int]:
        """找到x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        root = self.left[self._find(x + 1)]
        return root - 1 if root > 0 else None

    def next(self, x: int) -> Optional[int]:
        """x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        root = self.right[self._find(x)]
        return root - 1 if root < self._n + 1 else None

    def erase(self, start: int, end=-1) -> None:
        """删除[left, right)区间内的元素.
        0<=left<=right<=n.
        """
        if end == -1:
            end = start + 1
        if start >= end:
            return
        while True:
            m = self.right[self._find(start)]
            if m > end:
                break
            self._union(start, m)

    def has(self, x: int) -> bool:
        return self._find(x + 1) == x + 1

    def _find(self, x: int) -> int:
        if self._data[x] < 0:
            return x
        self._data[x] = self._find(self._data[x])
        return self._data[x]

    def _union(self, x: int, y: int) -> bool:
        rootX = self._find(x)
        rootY = self._find(y)
        if rootX == rootY:
            return False
        if self._data[rootX] > self._data[rootY]:
            rootX, rootY = rootY, rootX
        self._data[rootX] += self._data[rootY]
        self._data[rootY] = rootX
        if self.left[rootY] < self.left[rootX]:
            self.left[rootX] = self.left[rootY]
        if self.right[rootY] > self.right[rootX]:
            self.right[rootX] = self.right[rootY]
        return True

    def __contains__(self, x: int) -> bool:
        return self.has(x)

    def __repr__(self) -> str:
        res = [i for i in range(self._n) if self.has(i)]
        return f"Finder{res}"


if __name__ == "__main__":

    class Solution:
        # 3639. 变为活跃状态的最小时间
        # https://leetcode.cn/problems/minimum-time-to-activate-string/description/
        def minTime(self, s: str, order: List[int], k: int) -> int:
            n = len(s)
            count = n * (n + 1) // 2
            if count < k:
                return -1

            finder = Finder(n)
            for t in range(n - 1, -1, -1):
                i = order[t]
                l, r = finder.prev(i - 1), finder.next(i + 1)
                if l is None:
                    l = -1
                if r is None:
                    r = n
                count -= (i - l) * (r - i)
                if count < k:
                    return t
                finder.erase(i)
            raise RuntimeError("Should not reach here.")

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
        F = Finder(n)
        ok = [True] * n
        for _ in range(100):
            left, right = randint(0, n), randint(0, n)
            F.erase(left, right)
            erase(left, right)
            for i in range(n):
                assert F.prev(i) == pre(i), (i, F.prev(i), pre(i))
                assert F.next(i) == nxt(i), (i, F.next(i), nxt(i))
    print("Done!")
