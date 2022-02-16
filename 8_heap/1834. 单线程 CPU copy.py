import heapq
from typing import List


class Solution:
    def getOrder(self, tasks: List[List[int]]) -> List[int]:
        sortedTasks = [(task[0], i, task[1]) for i, task in enumerate(tasks)]
        sortedTasks.sort()
        pq = []
        time = 0
        res = []
        pos = 0
        for _ in sortedTasks:
            if not pq:
                time = max(time, sortedTasks[pos][0])
            while pos < len(sortedTasks) and sortedTasks[pos][0] <= time:
                heapq.heappush(pq, (sortedTasks[pos][2], sortedTasks[pos][1]))
                pos += 1
            d, j = heapq.heappop(pq)
            time += d
            res.append(j)
        return res
        
