# MK平均值
# !滑动窗口topK之和 (这里topK是最小值)

from typing import List
from sortedcontainers import SortedList


def windowTopKSum(nums: List[int], windowSize: int, k: int) -> List[int]:
    def add(x: int) -> None:
        nonlocal topKSum
        pos = sl.bisect_left(x)
        if pos < k:
            topKSum += x
            if k - 1 < len(sl):
                topKSum -= sl[k - 1]  # type: ignore
        sl.add(x)

    def remove(x: int) -> None:
        nonlocal topKSum
        pos = sl.bisect_left(x)
        if pos < k:
            topKSum -= x
            if k < len(sl):
                topKSum += sl[k]  # type: ignore
        sl.remove(x)

    def query() -> int:
        return topKSum

    n = len(nums)
    sl = SortedList()
    res, topKSum = [], 0
    for right in range(n):
        add(nums[right])
        if right >= windowSize:
            remove(nums[right - windowSize])
        if right >= windowSize - 1:
            res.append(query())
    return res


assert windowTopKSum([3, 1, 4, 1, 5, 9], 4, 3) == [5, 6, 10]
assert windowTopKSum([12, 2, 17, 11, 19, 8, 4, 3, 6, 20], 6, 3) == [21, 14, 15, 13, 13]


class TopkSum:
    __slots__ = ("_sl", "_k", "_topKSum")

    def __init__(self, k: int, min: bool) -> None:
        self._sl = SortedList() if min else SortedList(key=lambda x: -x)
        self._k = k
        self._topKSum = 0

    def add(self, x: int) -> None:
        pos = self._sl.bisect_left(x)
        if pos < self._k:
            self._topKSum += x
            if self._k - 1 < len(self._sl):
                self._topKSum -= self._sl[self._k - 1]  # type: ignore
        self._sl.add(x)

    def remove(self, x: int) -> None:
        pos = self._sl.bisect_left(x)
        if pos < self._k:
            self._topKSum -= x
            if self._k < len(self._sl):
                self._topKSum += self._sl[self._k]  # type: ignore
        self._sl.remove(x)

    def discard(self, x: int) -> None:
        if x in self._sl:
            self.remove(x)

    def query(self) -> int:
        return self._topKSum

    def __len__(self) -> int:
        return len(self._sl)


topKSum = TopkSum(3, min=True)
assert topKSum.query() == 0
topKSum.add(3)
assert topKSum.query() == 3
topKSum.add(1)
assert topKSum.query() == 4
topKSum.add(4)
assert topKSum.query() == 8
topKSum.add(1)
assert topKSum.query() == 5
topKSum.remove(3)
assert topKSum.query() == 6
