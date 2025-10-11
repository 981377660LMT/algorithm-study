# 3708. 最长斐波那契子数组
# https://leetcode.cn/problems/longest-fibonacci-subarray/description/
# 给你一个由 正 整数组成的数组 nums。
# 斐波那契 数组是一个连续序列，其中第三项及其后的每一项都等于这一项前面两项之和。
# 返回 nums 中最长 斐波那契 子数组的长度。
# 注意: 长度为 1 或 2 的子数组总是 斐波那契 的。
# 子数组 是数组中 非空 的连续元素序列。
#
# !方法1:双指针
# !方法2:斐波那契数列每两个元素就会翻倍，呈指数级增长。因此斐波那契子数组的长度是O(logX)的，可以暴力检查.

from typing import List


class Solution:
    def longestSubarray(self, nums: List[int]) -> int:
        n = len(nums)
        res = 2
        left = 0
        for right in range(2, n):
            if nums[right] != nums[right - 1] + nums[right - 2]:
                left = right - 1
            res = max(res, right - left + 1)
        return res
