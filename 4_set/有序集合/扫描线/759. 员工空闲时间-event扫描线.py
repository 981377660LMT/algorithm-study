# Definition for an Interval.
from collections import defaultdict
from typing import List, Optional
from sortedcontainers import SortedDict


class Interval:
    def __init__(self, start: Optional[int] = None, end: Optional[int] = None):
        self.start = start
        self.end = end


# 每个员工都有一个非重叠的时间段  Intervals 列表，这些时间段已经排好序。
# 返回表示 所有 员工的 共同，正数长度的空闲时间 的有限时间段的列表，同样需要排好序。


class Solution:
    def employeeFreeTime(self, schedule: List[List[Interval]]) -> List[Interval]:
        events = []
        for S in schedule:
            for interval in S:
                # 先进后出
                events.append((interval.start, 0))
                events.append((interval.end, 1))
        events.sort()

        res = []
        pre, preSum = -1, 0
        for key, flag in events:
            if preSum == 0 and pre != -1:
                res.append(Interval(pre, key))
            preSum += 1 if flag == 0 else -1
            pre = key

        return res


# 输入：schedule = [[[1,2],[5,6]],[[1,3]],[[4,10]]]
# 输出：[[3,4]]
# 解释：
# 共有 3 个员工，并且所有共同的
# 空间时间段是 [-inf, 1], [3, 4], [10, inf]。
# 我们去除所有包含 inf 的时间段，因为它们不是有限的时间段。

