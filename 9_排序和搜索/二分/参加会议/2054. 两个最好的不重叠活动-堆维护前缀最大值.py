from typing import List, Tuple
from heapq import heappop, heappush

# 你 最多 可以参加 两个时间不重叠 活动，使得它们的价值之和 最大 。
# 请你返回价值之和的 最大值 。
# 如果你参加一个活动，且结束时间为 t ，那么下一个活动必须在 t + 1 或之后的时间开始。
# 你 最多 可以参加 `两个`时间不重叠 活动，使得它们的价值之和 最大 。

# 13_简单游戏-前缀和+pq


class Solution:
    def maxTwoEvents(self, events: List[List[int]]) -> int:
        events.sort()
        pq: List[Tuple[int, int]] = []
        res, pre_max = 0, 0
        for start, end, val in events:
            heappush(pq, (end, val))
            while pq and pq[0][0] < start:
                _, pre_val = heappop(pq)
                pre_max = max(pre_max, pre_val)
            res = max(res, pre_max + val)
        return res


print(Solution().maxTwoEvents([[1, 3, 2], [4, 5, 2], [2, 4, 3]]))
# 如果要参加3个，那么pre_max就需要2个数记录
# 如果没有参加限制，则转化为出租车问题dp+二分 参见 1235. 规划兼职工作
