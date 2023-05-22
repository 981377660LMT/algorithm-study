# 区间合并


from typing import List


def mergeIntervals(intervals: List[List[int]]) -> List[List[int]]:
    """合并所有重叠的区间，并返回 一个不重叠的区间数组"""
    if not intervals:
        return []

    intervals.sort()
    res = [intervals[0]]
    for s, e in intervals[1:]:
        if s <= res[-1][1]:
            res[-1][1] = max(res[-1][1], e)
        else:
            res.append([s, e])

    return res


print(mergeIntervals([[1, 3], [2, 6], [8, 10], [15, 18]]))
