from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def stored_energy(self, store_limit: int, power: List[int], supply: List[List[int]]) -> int:
        n = len(power)

        change = defaultdict(tuple)
        for time, lower, upper in supply:
            change[time] = (lower, upper)

        remain = 0
        min_, max_ = 0, 0
        for i in range(n):
            if i in change:
                min_ = change[i][0]
                max_ = change[i][1]
            cur = power[i]
            if cur > max_:
                remain = min(remain + cur - max_, store_limit)
            elif cur < min_:
                remain = max(remain + cur - min_, 0)

        return remain
