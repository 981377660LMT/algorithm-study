from typing import List

# 给你一个整数数组 nums 和一个整数 k ，请你返回子数组内所有元素的乘积严格小于 k 的连续子数组的数目。
# 1 <= nums.length <= 3 * 104
# 1 <= nums[i] <= 1000
class Solution:
    def numSubarrayProductLessThanK(self, nums: List[int], k: int) -> int:
        """滑动具有单调性，所以滑窗"""
        left, res, curProduct = 0, 0, 1
        for right, num in enumerate(nums):
            curProduct *= num
            while left <= right and curProduct >= k:
                curProduct /= nums[left]
                left += 1
            res += right - left + 1
        return res


print(Solution().numSubarrayProductLessThanK([10, 5, 2, 6], 100))
print(Solution().numSubarrayProductLessThanK(nums=[1, 2, 3], k=0))
