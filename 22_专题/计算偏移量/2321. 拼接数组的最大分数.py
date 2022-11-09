from typing import List


# 可以交换一段子数组[left,right]
# !问能取到sum(nums1)和sum(nums2)的最大子数组和
# 公式变形:
# !sum(nums1) = sum(preNums1) + diff(nums2-nums1 这一段子数组)
# 所以是求两个数组的diff的最大子数组


def helper(nums: List[int]) -> int:
    res = 0
    dp = 0
    for num in nums:
        dp = max(dp + num, num)
        res = max(res, dp)
    return res


class Solution:
    def maximumsSplicedArray(self, nums1: List[int], nums2: List[int]) -> int:
        """求出两个数组的diff 那么就是要求diff的最大子数组和"""

        def cal(A: List[int], B: List[int]) -> int:
            diff = [b - a for a, b in zip(A, B)]
            return sum(A) + helper(diff)

        return max(cal(nums1, nums2), cal(nums2, nums1))
