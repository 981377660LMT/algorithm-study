from typing import List, Tuple
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def canReceiveAllSignals(self, intervals: List[List[int]]) -> bool:
        """区间不重叠"""
        if not intervals:
            return True
        intervals.sort()
        preEnd = intervals[0][1]
        for curStart, curEnd in intervals[1:]:
            if curStart < preEnd:
                return False
            preEnd = max(preEnd, curEnd)
        return True
