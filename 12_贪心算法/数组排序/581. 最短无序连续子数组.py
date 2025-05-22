# 581. 最短无序连续子数组
# 找到最短的子数组，排序后
# 给你一个整数数组 nums ，你需要找出一个 连续子数组 ，如果对这个子数组进行升序排序，那么整个数组都会变为升序排序。
# 请你找出符合题意的 最短 子数组，并输出它的长度。
#
# !找到最右侧的，比前缀最大值严格小的下标right
# !找到最左侧的，比后缀最小值严格小的下标left

from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


def min2(a: int, b: int) -> int:
    return a if a < b else b


class Solution:
    def findUnsortedSubarray(self, nums: List[int]) -> int:
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
