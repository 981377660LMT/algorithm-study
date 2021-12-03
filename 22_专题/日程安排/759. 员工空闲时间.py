# Definition for an Interval.
from typing import List
from sortedcontainers import SortedDict


class Interval:
    def __init__(self, start: int = None, end: int = None):
        self.start = start
        self.end = end


# 每个员工都有一个非重叠的时间段  Intervals 列表，这些时间段已经排好序。
# 返回表示 所有 员工的 共同，正数长度的空闲时间 的有限时间段的列表，同样需要排好序。

# 732. 我的日程安排表3.ts
class Solution:
    def employeeFreeTime(self, schedule: List[List[Interval]]) -> List[Interval]:
        dic = SortedDict()
        for sche in schedule:
            for interval in sche:
                start = interval.start
                end = interval.end
                if start not in dic:
                    dic[start] = 0
                dic[start] += 1
                if end not in dic:
                    dic[end] = 0
                dic[end] -= 1

        res = []
        pre_time = -1
        cur_worker = 0

        for cur_time, diff in dic.items():
            cur_worker += diff

            # 空闲时记录最左端点
            if cur_worker == 0:
                if pre_time == -1:
                    pre_time = cur_time
            # 上班时，结束空闲interval
            else:
                if pre_time != -1:
                    res.append(Interval(pre_time, cur_time))
                    pre_time = -1

        return res


# 输入：schedule = [[[1,2],[5,6]],[[1,3]],[[4,10]]]
# 输出：[[3,4]]
# 解释：
# 共有 3 个员工，并且所有共同的
# 空间时间段是 [-inf, 1], [3, 4], [10, inf]。
# 我们去除所有包含 inf 的时间段，因为它们不是有限的时间段。

