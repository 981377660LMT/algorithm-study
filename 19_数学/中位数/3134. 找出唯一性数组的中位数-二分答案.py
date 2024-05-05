# 3134. 找出唯一性数组的中位数-二分答案
# https://leetcode.cn/problems/find-the-median-of-the-uniqueness-array/
#
# 数组nums的唯一性数组是一个按元素从小到大排序的数组，
# 包含了nums的所有非空子数组中不同元素的个数。
# 返回 nums 唯一性数组 的 中位数 。
# !注意，数组的 中位数 定义为有序数组的中间元素。如果有两个中间元素，则取值较小的那个。
#
#
# !求中位数，考虑二分答案 =>
# !中位数为mid时，子数组的distinct元素个数小于等于mid的子数组个数必须等于k个.

from collections import defaultdict
from typing import List


class Solution:
    def medianOfUniquenessArray(self, nums: List[int]) -> int:
        def countNGT(mid: int) -> int:
            """有多少个子数组的distinct元素个数小于等于mid."""
            res, left, n = 0, 0, len(nums)
            counter = defaultdict(int)
            for right in range(n):
                counter[nums[right]] += 1
                while left <= right and len(counter) > mid:
                    counter[nums[left]] -= 1
                    if counter[nums[left]] == 0:
                        del counter[nums[left]]
                    left += 1
                res += right - left + 1
            return res

        n = len(nums)
        left, right = 0, n
        k = (n * (n + 1) // 2 + 1) // 2  # 二分查找第k小的数(k从1开始)
        while left <= right:
            mid = (left + right) // 2
            if countNGT(mid) < k:
                left = mid + 1
            else:
                right = mid - 1
        return left
