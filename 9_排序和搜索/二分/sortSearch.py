from bisect import bisect_left
from typing import Callable, Tuple


def sortSearch(start: int, end: int, f: Callable[[int], bool]) -> Tuple[int, bool]:
    """在区间 `[start, end)` 中查找使 `f(i)` 为 `true` 的最小值 `i`."""
    i = start
    j = end
    while i < j:
        h = (i + j) >> 1
        if not f(h):
            i = h + 1
        else:
            j = h
    return i, i < end and f(i)


def sortSearch2(start: int, end: int, f: Callable[[int], bool]) -> int:
    return bisect_left(range(start, end), True, key=f)


def sortSearchInts(arr: list[int], target: int) -> int:
    """bisect.bisect_left"""
    return sortSearch(0, len(arr), lambda i: arr[i] >= target)[0]
