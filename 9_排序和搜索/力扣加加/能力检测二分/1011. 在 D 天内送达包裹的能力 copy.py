from typing import List


class Solution:
    def shipWithinDays(self, weights: List[int], days: int) -> int:
        def check(mid: int) -> bool:
            res = 1
            curSum = 0
            for num in weights:
                if num > mid:
                    return False
                if curSum + num > mid:
                    res += 1
                    curSum = num
                else:
                    curSum += num
            return res <= days

        left, right = 1, int(1e19)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left

