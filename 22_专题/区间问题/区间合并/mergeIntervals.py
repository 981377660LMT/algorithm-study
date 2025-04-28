from typing import Generator, List, Tuple


def max2(a: int, b: int) -> int:
    return a if a > b else b


def mergeIntervals(intervals: List[List[int]]) -> Generator[Tuple[int, int], None, None]:
    """合并所有重叠的区间，并返回一个不重叠的区间数组.

    >>> list(mergeIntervals([[1, 2], [2, 4], [5, 6]]))
    [(1, 4), (5, 6)]
    """

    if not intervals:
        return
    order = sorted(range(len(intervals)), key=lambda i: intervals[i][0])
    preL, preR = intervals[order[0]]
    for i in order[1:]:
        curL, curR = intervals[i]
        if curL <= preR:
            preR = max2(preR, curR)
        else:
            yield (preL, preR)
            preL, preR = curL, curR
    yield (preL, preR)


if __name__ == "__main__":
    import doctest

    doctest.testmod()
