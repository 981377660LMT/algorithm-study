from bisect import bisect_left, bisect_right
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def explorationSupply(self, station: List[int], pos: List[int]) -> List[int]:
        res = []
        for p in pos:
            index = bisect_right(station, p) - 1
            cand = index
            dist = abs(station[index] - p)
            if index + 1 < len(station) and abs(station[index + 1] - p) < dist:
                cand = index + 1
            res.append(cand)
        return res


print(Solution().explorationSupply([5, 9, 10, 12, 15], [8, 9, 4, 16, 17]))
