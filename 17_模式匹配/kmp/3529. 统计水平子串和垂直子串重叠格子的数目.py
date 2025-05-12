from itertools import accumulate
from typing import List

from kmp import indexOfAll


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


class Solution:
    def countCells(self, grid: List[List[str]], pattern: str) -> int:
        row, col = len(grid), len(grid[0])

        rowText = [c for row in grid for c in row]
        rowIndexes = indexOfAll(rowText, pattern)
        diff1 = DiffArray(row * col)
        for start in rowIndexes:
            diff1.add(start, start + len(pattern), 1)

        colText = [c for col in zip(*grid) for c in col]
        colIndexes = indexOfAll(colText, pattern)
        diff2 = DiffArray(row * col)
        for start in colIndexes:
            diff2.add(start, start + len(pattern), 1)

        arr1 = diff1.getAll()
        arr2 = diff2.getAll()
        res = 0
        for r in range(row):
            for c in range(col):
                if arr1[r * col + c] > 0 and arr2[c * row + r] > 0:
                    res += 1
        return res
