# 3323. 通过插入区间最小化连通组
# https://leetcode.cn/problems/minimize-connected-groups-by-inserting-interval/


from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


def mergeIntervals(intervals: List[List[int]]) -> List[List[int]]:
    """合并所有重叠的区间，并返回一个不重叠的区间数组.

    >>> mergeIntervals([[1, 2], [2, 4], [5, 6]])
    [[1, 4], [5, 6]]
    """
    if not intervals:
        return []
    intervals = sorted(intervals, key=lambda x: x[0])
    res = []
    for interval in intervals:
        if not res or res[-1][1] < interval[0]:
            res.append(interval)
        else:
            res[-1][1] = max2(res[-1][1], interval[1])
    return res


class Solution:
    def minConnectedGroups(self, intervals: List[List[int]], k: int) -> int:
        intervals = mergeIntervals(intervals)
        maxOverlap, left, n = 0, 0, len(intervals)
        for right in range(n):
            while left <= right and intervals[right][0] - intervals[left][1] > k:
                left += 1
            maxOverlap = max2(maxOverlap, right - left)
        return len(intervals) - maxOverlap
