# 非负整数数组 nums 。如果存在一个数 x ，
# 使得 nums 中恰好有 x 个元素 大于或者等于 x ，
# 那么就称 nums 是一个 特殊数组 ，而 x 是该数组的 特征值 。
from typing import List


# 1 <= nums.length <= 10^6
class Solution:
    def specialArray(self, nums: List[int]) -> int:
        left, right = 0, len(nums)
        while left <= right:
            mid = (left + right) >> 1
            count = 0
            for num in nums:
                if num >= mid:
                    count += 1
            if count == mid:
                return mid
            elif count > mid:
                left = mid + 1
            else:
                right = mid - 1
        return -1


# 输入：nums = [3,5]
# 输出：2
# 解释：有 2 个元素（3 和 5）大于或等于 2 。
