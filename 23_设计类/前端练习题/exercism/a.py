from typing import List
from bisect import bisect_right
from collections import defaultdict
from itertools import accumulate


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

    def sumAll(self) -> int:
        """求和."""
        self.build()
        return self._preSum[-1]


# 维护增量
class Solution:
    def maxArea(self, height: int, positions: List[int], directions: str) -> int:
        ...


# height = 5, positions = [2,5], directions = "UD"

print(Solution().maxArea(5, [2, 5], "UD") == 7)
