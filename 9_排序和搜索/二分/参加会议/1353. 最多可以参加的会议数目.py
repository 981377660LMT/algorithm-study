from typing import List
from heapq import heappop, heappush


# 每个会议你至少参加一天
# 一天只能参加一个会议。

# 思路：排序；贪心，始终参加结束时间最早的会议
class Solution:
    def maxEvents(self, events: List[List[int]]) -> int:
        events.sort(key=lambda x: x[0])
        ei, res, pq = 0, 0, []
        for d in range(int(1e5) + 10):
            while ei < len(events) and events[ei][0] == d:
                heappush(pq, events[ei][1])
                ei += 1
            while pq and pq[0] < d:
                heappop(pq)
            if pq:
                heappop(pq)
                res += 1
        return res


print(Solution().maxEvents([[1, 4], [4, 4], [2, 2], [3, 4], [1, 1]]))
