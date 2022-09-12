# 1. 左端点从小到大排序
# 2. 遍历区间，判断能否判断放入现有组中

# 给定 N 个闭区间 [ai,bi]，
# 请你将这些区间分成若干组，
# 使得每组内部的区间两两之间（包括端点）没有交集，
# 并使得组数尽可能小。

# 输出最小组数。
# !n<=1e4
# !start<end<1e9
#######################################################################
# 解答：
# !区间按照左端点排序
# 小根堆存放所有组的右端点值，堆顶存放最小的右端点值
# 如果当前区间左端点大于堆顶元素，说明可以加入堆顶元素所在组，右端点入堆
# 如果当前区间左端点小于等于堆顶元素，说明当前区间与堆里面的区间重叠


# !也可以差分做:
# 会议室问题
# https://leetcode.cn/problems/meeting-rooms-ii/
from collections import defaultdict
from heapq import heappop, heappush
from itertools import accumulate
from typing import List


class Solution:
    def minGroups(self, intervals: List[List[int]]) -> int:
        intervals.sort()
        pq = []
        for start, end in intervals:
            if pq and start > pq[0]:
                heappop(pq)
            heappush(pq, end)  # 更新分区的末尾
        return len(pq)

    def minGroups2(self, intervals: List[List[int]]) -> int:
        """会议室系列 差分更新"""
        diff = defaultdict(int)
        for start, end in intervals:
            diff[start] += 1
            diff[end + 1] -= 1
        nums = [diff[key] for key in sorted(diff)]
        return max(accumulate(nums))
