# https://nyaannyaan.github.io/library/data-structure/segment-set.hpp
# 管理区间的数据结构


from typing import Tuple, Union
from sortedcontainers import SortedList

INF = int(1e18)


class SegmentSet:
    __slots__ = ("_st",)

    def __init__(self):
        self._st = SortedList()

    def insert(self, left: int, right: int) -> None:
        """插入闭区间[left, right]."""
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

    def erase(self, left: int, right: int) -> None:
        """删除闭区间[left, right]."""
        it1 = self._st.bisect_right((left, INF))
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

    def next(self, x: int) -> int:
        """返回第一个大于等于x的区间起点.如果不存在,返回INF."""
        it = self._st.bisect_left((x, -INF))
        if it == len(self._st):
            return INF
        res = self._st[it][0]
        if x > res:
            return x
        return res

    def __contains__(self, arg: Union[int, Tuple[int, int]]) -> bool:
        if isinstance(arg, int):
            it = self._st.bisect_right((arg, INF))
            return it != 0 and self._st[it - 1][1] >= arg
        left, right = arg
        assert left <= right, "left must be less than or equal to right"
        it1 = self._st.bisect_right((left, INF))
        if it1 == 0:
            return False
        it2 = self._st.bisect_right((right, INF))
        if it1 != it2:
            return False
        return self._st[it1 - 1][1] >= right

    def __getitem__(self, index: int) -> Tuple[int, int]:
        return self._st[index]

    def __repr__(self) -> str:
        return repr(self._st)

    def __len__(self) -> int:
        return len(self._st)


if __name__ == "__main__":
    s = SegmentSet()
    s.insert(1, 3)
    s.insert(2, 4)
    s.insert(5, 6)
    s.insert(7, 8)
    s.insert(6, 7)
    s.insert(0, 9)
    s.insert(0, 10)
    s.erase(0, 2)
    s.erase(0, 9)
    s.insert(1, 3)
    s.insert(2, 4)

    print(s)
    print((9, 11) in s)
