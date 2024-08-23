from typing import Callable, List
from itertools import accumulate


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class PresumBisector:
    """带有二分的前缀和，要求元素为非负数."""

    __slots__ = "_n", "_presum"

    def __init__(self, nums: List[int]):
        self._n = len(nums)
        self._presum = list(accumulate(nums, initial=0))

    def query(self, start: int, end: int) -> int:
        start = max2(0, start)
        end = min2(self._n, end)
        if start >= end:
            return 0
        return self._presum[end] - self._presum[start]

    def maxRight(self, left: int, check: Callable[[int, int], bool]) -> int:
        if left >= self._n:
            return self._n
        ok, ng = left, self._n + 1
        while ok + 1 < ng:
            mid = (ok + ng) >> 1
            if check(self.query(left, mid), mid):
                ok = mid
            else:
                ng = mid
        return ok

    def minLeft(self, right: int, check: Callable[[int, int], bool]) -> int:
        if right <= 0:
            return 0
        ok, ng = right, -1
        while ng + 1 < ok:
            mid = (ok + ng) >> 1
            if check(self.query(mid, right), mid):
                ok = mid
            else:
                ng = mid
        return ok
