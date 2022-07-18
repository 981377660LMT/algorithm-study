from heapq import nlargest
import re
from typing import List, Tuple, Optional
from collections import defaultdict, Counter, deque

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def maximumSum(self, nums: List[int]) -> int:
        adjMap = defaultdict(list)
        for num in nums:
            sum_ = sum(map(int, str(num)))
            adjMap[sum_].append(num)

        res = -1
        for v in adjMap.values():
            if len(v) <= 1:
                continue
            cand = sum(nlargest(2, v))
            if cand > res:
                res = cand

        return res


print(Solution().maximumSum(nums=[18, 43, 36, 13, 7]))
