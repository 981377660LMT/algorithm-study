# 更快的O(n)算法，C++ 4ms 100%
# https://leetcode.cn/problems/koko-eating-bananas/solutions/1539347/by-hqztrue-3gs3/


from typing import List


class Solution:
    def minEatingSpeed(self, piles: List[int], h: int) -> int:
        """O(nlogU)."""

        def check(mid: int) -> bool:
            res = 0
            for num in piles:
                res += (num + mid - 1) // mid
            return res <= h

        left, right = 1, max(piles)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
