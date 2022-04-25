from math import ceil
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)


# LCP 55. 采集果实


class Solution:
    def getMinimumTime(self, time: List[int], fruits: List[List[int]], limit: int) -> int:
        res = 0
        for t, v in fruits:
            res += ceil(v / limit) * time[t]
        return res
