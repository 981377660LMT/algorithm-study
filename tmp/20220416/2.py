from math import ceil
from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def waysToBuyPensPencils(self, total: int, cost1: int, cost2: int) -> int:
        count1 = 0
        res = 0
        while count1 * cost1 <= total:
            remain = total - count1 * cost1
            res += remain // cost2 + 1
            count1 += 1
        return res

