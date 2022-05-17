from bisect import bisect_left
from typing import List


class Solution:
    def findRightInterval(self, intervals: List[List[int]]) -> List[int]:
        n = len(intervals)
        res = [-1] * n

        starts = sorted([s, i] for i, (s, _) in enumerate(intervals))
        for i, (_, e) in enumerate(intervals):
            # pos = bisect_left(starts, (e, -int(1e20)))
            pos = bisect_left(starts, e, key=lambda x: x[0])
            if pos < n:
                res[i] = starts[pos][1]
        return res

