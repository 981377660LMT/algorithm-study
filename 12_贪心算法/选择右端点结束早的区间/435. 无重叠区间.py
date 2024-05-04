# 435. 无重叠区间
# https://leetcode.cn/problems/non-overlapping-intervals/description/
# 给定一个区间的集合 intervals ，其中 intervals[i] = [starti, endi] 。
# 返回 需要移除区间的最小数量，使剩余区间互不重叠 。

from typing import List


INF = int(1e18)


class Solution:
    def eraseOverlapIntervals(self, intervals: List[List[int]]) -> int:
        intervals.sort(key=lambda x: x[1])
        preEnd = -INF
        minOverlap = 0
        for curStart, curEnd in intervals:
            if curStart < preEnd:
                minOverlap += 1
            else:
                preEnd = curEnd
        return minOverlap
