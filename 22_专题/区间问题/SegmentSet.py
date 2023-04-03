# 管理区间的数据结构
# 注意:
# 1.所有区间都是闭区间 例如 [1,1] 表示 长为1的区间,起点为1
# !2.有交集的区间会被合并,例如 [1,2]和[2,3]会被合并为[1,3]


from typing import Optional, Tuple, Union
from sortedcontainers import SortedList

INF = int(1e18)


class SegmentSet:
    __slots__ = ("_st", "count")

    def __init__(self):
        self.count = 0  # 区间元素的个数
        self._st = SortedList()

    def insert(self, left: int, right: int) -> None:
        """插入闭区间[left, right]."""
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
            removed = sum(right - left + 1 for left, right in self._st[it1:it2])
            del self._st[it1:it2]
            self.count -= removed
        self._st.add((left, right))
        self.count += right - left + 1

    def erase(self, left: int, right: int) -> None:
        """删除闭区间[left, right]."""
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
        removed = sum(right - left + 1 for left, right in self._st[it1:it2])
        del self._st[it1:it2]
        self.count -= removed
        if nl < left:
            self._st.add((nl, left))
            self.count += left - nl + 1
        if right < nr:
            self._st.add((right, nr))
            self.count += nr - right + 1

    def nextStart(self, x: int) -> Optional[int]:
        """返回第一个大于等于x的`区间起点`.如果不存在,返回None."""
        it = self._st.bisect_left((x, -INF))
        if it == len(self._st):
            return
        res = self._st[it][0]
        if x > res:
            return x
        return res

    def ceiling(self, x: int) -> Optional[int]:
        """返回区间内第一个大于等于x的元素.如果不存在,返回None."""
        pos = self._st.bisect_right((x, INF))
        if pos != 0 and self._st[pos - 1][1] >= x:
            return x
        if pos != len(self._st):
            return self._st[pos][0]

    def floor(self, x: int) -> Optional[int]:
        """返回区间内最后一个小于等于x的元素.如果不存在,返回None."""
        pos = self._st.bisect_right((x, INF))
        if pos == 0:
            return None
        if self._st[pos - 1][1] >= x:
            return x
        return self._st[pos - 1][1]

    def get(self, x: int) -> Optional[Tuple[int, int]]:
        """返回包含x的区间.如果不存在,返回None."""
        pos = self._st.bisect_right((x, INF))
        if pos == 0 or self._st[pos - 1][1] < x:
            return None
        return self._st[pos - 1]

    def __contains__(self, arg: Union[int, Tuple[int, int]]) -> bool:
        if isinstance(arg, int):
            it = self._st.bisect_right((arg, INF))
            return it != 0 and self._st[it - 1][1] >= arg
        left, right = arg
        if left > right:
            return True
        it1 = self._st.bisect_right((left, INF))
        if it1 == 0:
            return False
        it2 = self._st.bisect_right((right, INF))
        if it1 != it2:
            return False
        return self._st[it1 - 1][1] >= right

    def __getitem__(self, index: int) -> Tuple[int, int]:
        return self._st[index]

    def __iter__(self):
        return iter(self._st)

    def __repr__(self) -> str:
        sb = []
        for left, right in self._st:
            sb.append(f"({left}, {right})")
        return f"SegmentSet([{', '.join(sb)}])"

    def __len__(self) -> int:
        return len(self._st)


if __name__ == "__main__":
    ss = SegmentSet()
    ss.insert(1, 3)
    ss.insert(2, 4)
    ss.insert(5, 6)
    assert ss.nextStart(1) == 1
    assert (1, 4) in ss
    assert 7 not in ss
    assert ss.count == sum(right - left + 1 for left, right in ss)
    ss.erase(2, 3)
    assert len(ss) == 3
    assert ss.get(5) == (5, 6)

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
        F = SegmentSet()
        for i in range(n):
            F.insert(i, i)
        ok = [True] * n
        for _ in range(100):
            left, right = randint(0, n), randint(0, n)
            F.erase(left, right - 1)
            erase(left, right)
            for i in range(n):
                assert F.floor(i) == pre(i), (i, F.floor(i), pre(i))
                assert F.ceiling(i) == nxt(i), (i, F.ceiling(i), nxt(i))
    print("Done!")
