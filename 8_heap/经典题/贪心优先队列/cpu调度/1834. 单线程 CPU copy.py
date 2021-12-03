import heapq
from typing import List


class Solution:
    def getOrder(self, tasks: List[List[int]]) -> List[int]:
        res = []
        events = sorted([(t[0], t[1], i) for i, t in enumerate(tasks)])

        event_idx = 0
        pq = []
        time = 0
        while len(res) < len(events):
            while (event_idx < len(events)) and (events[event_idx][0] <= time):
                heapq.heappush(
                    pq, (events[event_idx][1], events[event_idx][2])
                )  # (processing_time, original_index)
                event_idx += 1
            if pq:
                t_diff, original_index = heapq.heappop(pq)
                time += t_diff
                res.append(original_index)
            elif event_idx < len(events):
                time = events[event_idx][0]
        return res
