from typing import List


class Solution:
    def maxSubArray(self, nums: List[int]) -> int:
        """最大子数组和-取或全不取dp"""
        if len(nums) == 1:
            return nums[0]

        curMax, res = -int(1e20), nums[0]
        for num in nums:
            curMax = max(curMax + num, num)
            res = max(res, curMax)
        return res

    def maxSubArray2(self, nums: List[int]) -> int:
        """最大子数组和-前缀和"""
        if len(nums) == 1:
            return nums[0]

        curSum, res, preMin = 0, -int(1e20), 0
        for num in nums:
            curSum += num
            res = max(res, curSum - preMin)
            preMin = min(preMin, curSum)
        return res
