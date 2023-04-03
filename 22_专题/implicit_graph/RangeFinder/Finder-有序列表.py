# 寻找前驱后继/区间删除
# 注意:速度很慢


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

    def prev(self, x: int) -> Optional[int]:
        """找到x左侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        pos = self._st.bisect_right((x, INF))
        if pos == 0:
            return None
        if self._st[pos - 1][1] >= x:
            return x
        return self._st[pos - 1][1]

    def next(self, x: int) -> Optional[int]:
        """x右侧第一个未被访问过的位置(包含x).
        如果不存在,返回None.
        """
        pos = self._st.bisect_right((x, INF))
        if pos != 0 and self._st[pos - 1][1] >= x:
            return x
        if pos != len(self._st):
            return self._st[pos][0]

    def erase(self, left: int, right: int) -> None:
        """删除左闭右开区间[left, right)."""
        right -= 1
        if left > right:
            return
        it1 = self._st.bisect_left((left, -INF))
        it2 = self._st.bisect_right((right, INF))
        if it1 != 0 and left <= self._st[it1 - 1][1]:
            it1 -= 1
        if it1 == it2:
            return
        nl, nr = self._st[it1][0], self._st[it2 - 1][1]
        if left < nl:
            nl = left
        if right > nr:
            nr = right
        del self._st[it1:it2]
        if nl < left:
            self._st.add((nl, left))
        if right < nr:
            self._st.add((right, nr))

    def insert(self, left: int, right: int) -> None:
        """插入左闭右开区间[left, right)."""
        right -= 1
        if left > right:
            return
        it1 = self._st.bisect_right((left, INF))
        it2 = self._st.bisect_right((right, INF))
        if it1 != 0 and left <= self._st[it1 - 1][1]:
            it1 -= 1
        if it1 != it2:
            tmp1 = self._st[it1][0]
            if tmp1 < left:
                left = tmp1
            tmp2 = self._st[it2 - 1][1]
            if tmp2 > right:
                right = tmp2
            del self._st[it1:it2]
        self._st.add((left, right))

    def __repr__(self) -> str:
        sb = []
        for left, right in self._st:
            sb.append(f"({left}, {right})")
        return f"SegmentSet([{', '.join(sb)}])"


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
            F.insert(i, i + 1)
        ok = [True] * n
        for _ in range(100):
            left, right = randint(0, n), randint(0, n)
            F.erase(left, right)
            erase(left, right)
            for i in range(n):
                assert F.prev(i) == pre(i), (i, F.prev(i), pre(i))
                assert F.next(i) == nxt(i), (i, F.next(i), nxt(i))
    print("Done!")
