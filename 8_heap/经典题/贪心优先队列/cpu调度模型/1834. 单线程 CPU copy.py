# CPU调度-单线程

from heapq import heappop, heappush
from typing import List


class Solution:
    def getOrder(self, tasks: List[List[int]]) -> List[int]:
        res = []
        events = sorted([(t[0], t[1], i) for i, t in enumerate(tasks)])

        ei = 0
        pq = []
        time = 0
        while len(res) < len(events):
            while (ei < len(events)) and (events[ei][0] <= time):
                heappush(pq, (events[ei][1], events[ei][2]))  # (processing_time, original_index)
                ei += 1
            if pq:
                diff, preI = heappop(pq)
                time += diff
                res.append(preI)
            elif ei < len(events):
                time = events[ei][0]
        return res


# !不要维护time
