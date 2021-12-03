from typing import List
from heapq import heappop, heappush


# 每个会议你至少参加一天
# 一天只能参加一个会议。

# 思路：排序；贪心，始终参加结束时间最早的会议
class Solution:
    def maxEvents(self, events: List[List[int]]) -> int:
        events.sort()
        res = 0
        event_index = 0
        max_end = max([end for _, end in events])
        pq: List[int] = []

        for day in range(1, max_end + 1):
            # 当日开始的会议
            while event_index < len(events) and events[event_index][0] == day:
                heappush(pq, events[event_index][1])
                event_index += 1

            # 已经结束的会议
            while pq and pq[0] < day:
                heappop(pq)

            # 最早结束的会议
            if pq:
                res += 1
                heappop(pq)

        return res


print(Solution().maxEvents([[1, 4], [4, 4], [2, 2], [3, 4], [1, 1]]))
