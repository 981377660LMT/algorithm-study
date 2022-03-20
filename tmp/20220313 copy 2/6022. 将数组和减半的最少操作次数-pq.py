from heapq import heapify, heappop, heappush
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)


class Solution:
    def halveArray(self, nums: List[int]) -> int:
        target = sum(nums) / 2
        curSum = 0
        res = 0
        pq = [-n for n in nums]
        heapify(pq)
        # 这里最好epsilon 1e-6
        while curSum < target:
            top = -heappop(pq)
            curSum += top / 2
            heappush(pq, -top / 2)
            res += 1
        return res

