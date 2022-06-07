from math import ceil
from typing import List


class Solution:
    def minEatingSpeed(self, piles: List[int], h: int) -> int:
        def check(mid: int) -> bool:
            res = 0
            for num in piles:
                res += ceil(num / mid)
            return res <= h

        left, right = 1, max(piles) + 10
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left

