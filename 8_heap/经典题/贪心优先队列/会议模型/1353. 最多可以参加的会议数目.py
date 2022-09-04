from typing import List
from heapq import heappop, heappush


# 每个会议你至少参加一天
# 一天只能参加一个会议。

# !思路：排序；贪心，始终参加结束时间最早的会议

# 1 <= events.length <= 1e5
# events[i].length == 2
# 1 <= startDayi <= endDayi <= 1e5

# 在每一个时间点，我们首先将当前时间点开始的会议加入小根堆，
# 再把当前已经结束的会议移除出小根堆（因为已经无法参加了），
# 然后从剩下的会议中选择一个结束时间最早的去参加。


class Solution:
    def maxEvents(self, events: List[List[int]]) -> int:
        events.sort(key=lambda x: x[0])
        ei, res, pq = 0, 0, []
        for time in range(int(1e5) + 10):
            # 1.在每一个时间点，我们首先将当前时间点开始的会议加入小根堆，
            while ei < len(events) and events[ei][0] == time:
                heappush(pq, events[ei][1])
                ei += 1
            # 2.再把当前已经结束的会议移除出小根堆（因为已经无法参加了），
            while pq and pq[0] < time:
                heappop(pq)
            # 3.然后从剩下的会议中选择一个结束时间最早的去参加。
            if pq:
                heappop(pq)
                res += 1
        return res


print(Solution().maxEvents([[1, 4], [4, 4], [2, 2], [3, 4], [1, 1]]))
