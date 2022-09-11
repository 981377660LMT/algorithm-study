from heapq import heappop, heappush
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def minGroups(self, intervals: List[List[int]]) -> int:
        intervals.sort()

        pq = []
        for start, end in intervals:
            if pq and start > pq[0]:
                heappop(pq)
            heappush(pq, end)
        return len(pq)
