from bisect import bisect_right
from collections import defaultdict
from itertools import accumulate
from typing import List


class DiffArray:
    """差分维护区间修改，区间查询."""

    __slots__ = ("_diff", "_dirty")

    def __init__(self, n: int) -> None:
        self._diff = [0] * (n + 1)
        self._dirty = False

    def add(self, start: int, end: int, delta: int) -> None:
        """区间 `[start,end)` 加上 `delta`."""
        if start < 0:
            start = 0
        if end >= len(self._diff):
            end = len(self._diff) - 1
        if start >= end:
            return
        self._dirty = True
        self._diff[start] += delta
        self._diff[end] -= delta

    def build(self) -> None:
        if self._dirty:
            self._diff = list(accumulate(self._diff))
            self._dirty = False

    def get(self, pos: int) -> int:
        """查询下标 `pos` 处的值."""
        self.build()
        return self._diff[pos]

    def getAll(self) -> List[int]:
        self.build()
        return self._diff[:-1]


class DiffMap:
    """差分维护区间修改，单点查询."""

    __slots__ = ("_diff", "_preSum", "_sortedKeys", "_dirty")

    def __init__(self) -> None:
        self._diff = defaultdict(int)
        self._sortedKeys = []
        self._preSum = []
        self._dirty = False

    def add(self, start: int, end: int, delta: int) -> None:
        """区间 `[start,end)` 加上 `delta`."""
        if start >= end:
            return
        self._dirty = True
        self._diff[start] += delta
        self._diff[end] -= delta

    def build(self) -> None:
        if self._dirty:
            self._sortedKeys = sorted(self._diff)
            self._preSum = [0] + list(accumulate(self._diff[key] for key in self._sortedKeys))
            self._dirty = False

    def get(self, pos: int) -> int:
        """查询下标 `pos` 处的值."""
        self.build()
        return self._preSum[bisect_right(self._sortedKeys, pos)]


if __name__ == "__main__":
    # 2251. 花期内花的数目
    # https://leetcode.cn/problems/number-of-flowers-in-full-bloom/description/

    class Solution:
        def fullBloomFlowers(self, flowers: List[List[int]], people: List[int]) -> List[int]:
            diff = DiffMap()
            for left, right in flowers:
                diff.add(left, right + 1, 1)
            return [diff.get(p) for p in people]
