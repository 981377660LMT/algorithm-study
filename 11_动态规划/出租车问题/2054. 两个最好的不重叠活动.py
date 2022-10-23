from typing import List, Tuple
from heapq import heappop, heappush

# 你 最多 可以参加 两个时间不重叠 活动，使得它们的价值之和 最大 。
# 请你返回价值之和的 最大值 。
# 如果你参加一个活动，且结束时间为 t ，那么下一个活动必须在 t + 1 或之后的时间开始。
# 你 最多 可以参加 `两个`时间不重叠 活动，使得它们的价值之和 最大 。


# 关键思路：用堆维护之前的最大值，每次与当前相加
# !强化版:11_动态规划/出租车问题/1751. 最多可以参加的会议数目 II.py
class Solution:
    def maxTwoEvents(self, events: List[List[int]]) -> int:
        events.sort()
        pq: List[Tuple[int, int]] = []
        res, preMax = 0, 0
        for start, end, val in events:
            heappush(pq, (end, val))
            while pq and pq[0][0] < start:
                _, pre_val = heappop(pq)
                preMax = max(preMax, pre_val)
            res = max(res, preMax + val)
        return res


print(Solution().maxTwoEvents([[1, 3, 2], [4, 5, 2], [2, 4, 3]]))
# 如果要参加3个，那么pre_max就需要2个数记录
# 如果没有参加限制，则转化为出租车问题dp+二分 参见 1235. 规划兼职工作
# 需要按照结束顺序来dp
