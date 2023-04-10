# 寻找前驱后继/区间删除
# 注意:速度比较慢


from typing import Optional
from sortedcontainers import SortedList

INF = int(1e18)


class Finder:
    """利用SortedList寻找区间的某个位置左侧/右侧第一个未被访问过的位置.
    初始时,所有位置都未被访问过.
    """

    __slots__ = "_st"

    def __init__(self):
        self._st = SortedList()

    def insert(self, i: int) -> None:
        self._st.add(i)

    def erase(self, i: int) -> None:
        self._st.discard(i)

    def prev(self, x: int) -> Optional[int]:
        """找到x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        pos = self._st.bisect_right(x) - 1
        return self._st[pos] if pos >= 0 else None

    def next(self, x: int) -> Optional[int]:
        """x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        pos = self._st.bisect_left(x)
        return self._st[pos] if pos < len(self._st) else None

    def __contains__(self, x: int) -> bool:
        return x in self._st

    def __repr__(self) -> str:
        return repr(self._st)

    def __len__(self) -> int:
        return len(self._st)


if __name__ == "__main__":
    # 前驱后继
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
        F = Finder()
        for i in range(n):
            F.insert(i)
        ok = [True] * n
        for _ in range(100):
            i = randint(0, n - 1)
            F.erase(i)
            erase(i, i + 1)
            for i in range(n):
                assert F.prev(i) == pre(i), (i, F.prev(i), pre(i))
                assert F.next(i) == nxt(i), (i, F.next(i), nxt(i))
    print("Done!")
