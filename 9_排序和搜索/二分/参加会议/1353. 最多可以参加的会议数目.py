from typing import List
from heapq import heappop, heappush


# 你可以不完整参加会议
class Solution:
    def maxEvents(self, events: List[List[int]]) -> int:
        events.sort(key=lambda x: x[0])
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
