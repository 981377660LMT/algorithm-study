from itertools import pairwise
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def sumCounts(self, nums: List[int]) -> int:
        n = len(nums)
        res = 0
        subCount = 0
        last = defaultdict(lambda: -1)
        for i, num in enumerate(nums):
            subCount += i - last[num]
            res += subCount
            last[num] = i

        return res % MOD


# nums = [1,2,1]

print(Solution().sumCounts([1, 2, 1]))
