# 3555. 排序每个滑动窗口中最小的子数组
# https://leetcode.cn/problems/smallest-subarray-to-sort-in-every-sliding-window/description/
# 给定一个整数数组 nums 和一个整数 k。
# 对于每个长度为 k 的连续 子数组，确定必须排序的连续段的最小长度，以便整个窗口成为 非递减 的；如果窗口已经排序，则其所需长度为零。
# 返回一个长度为 n − k + 1 的数组，其中每个元素对应其窗口的答案。
#
# 1 <= nums.length <= 1000
# 1 <= k <= nums.length
# 1 <= nums[i] <= 1e6


from typing import List


def findUnsortedSubarray(nums: List[int]) -> int:
    """https://leetcode.cn/problems/shortest-unsorted-continuous-subarray/description/"""
    n = len(nums)
    if n <= 1:
        return 0

    preMax = nums[0]
    right = -1
    for i, v in enumerate(nums):
        if v >= preMax:
            preMax = v
        else:
            right = i

    sufMin = nums[-1]
    left = -1
    for i in range(n - 1, -1, -1):
        v = nums[i]
        if v <= sufMin:
            sufMin = v
        else:
            left = i

    return (right - left + 1) if (right != -1 and left != -1) else 0


class Solution:
    def minSubarraySort(self, nums: List[int], k: int) -> List[int]:
        return [findUnsortedSubarray(nums[l : l + k]) for l in range(len(nums) - k + 1)]
