# 3520. 逆序对计数的最小阈值
# https://leetcode.cn/problems/minimum-threshold-for-inversion-pairs-count/description
# RangeFreqDynamic

from typing import List

from sortedcontainers import SortedList


class Solution:
    def minThreshold(self, nums: List[int], k: int) -> int:
        def check(mid: int) -> bool:
            sl = SortedList()
            res = 0
            for num in nums:
                res += sl.bisect_left(num + mid + 1) - sl.bisect_left(num + 1)
                if res >= k:
                    return True
                sl.add(num)
            return res >= k

        left, right = 0, max(nums) - min(nums)
        if not check(right):
            return -1

        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left
