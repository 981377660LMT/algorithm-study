from typing import List

# nums全为正整数
class Solution:
    def kthSmallestSubarraySum(self, nums: List[int], k: int) -> None:
        def count(mid) -> int:
            """"区间和小于等于mid的子数组数"""

            res, curSum, left = 0, 0, 0
            for right in range(len(nums)):
                curSum += nums[right]
                while left < len(nums) and curSum > mid:
                    curSum -= nums[left]
                    left += 1
                res += right - left + 1
            return res
