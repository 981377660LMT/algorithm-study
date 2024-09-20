from typing import List


def mergeIntervals(intervals: List[List[int]]) -> List[List[int]]:
    """合并所有重叠的区间，并返回一个不重叠的区间数组.

    >>> mergeIntervals([[1, 2], [2, 4], [5, 6]])
    [[1, 4], [5, 6]]
    """
    if not intervals:
        return []

    intervals = intervals[:]
    intervals.sort(key=lambda x: x[0])
    res = []
    for interval in intervals:
        if not res or res[-1][1] < interval[0]:
            res.append(interval)
        else:
            res[-1][1] = max(res[-1][1], interval[1])

    return res


print(mergeIntervals([[1, 2], [2, 4], [5, 6]]))
