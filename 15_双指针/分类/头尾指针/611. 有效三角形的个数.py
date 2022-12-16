# 给定一个包含非负整数的数组 nums ，返回其中可以组成三角形三条边的三元组个数。
# n<=1000
# !排序+双指针定一移二

# !如果要求任意两边都不能相同,怎么做? => 二分查找

from bisect import bisect_left, bisect_right
from typing import List


class Solution:
    def triangleNumber(self, nums: List[int]) -> int:
        """排序+双指针O(n^2)"""
        n = len(nums)
        nums.sort()
        res = 0
        for i in range(n - 1, 1, -1):
            left, right = 0, i - 1
            while left < right:
                if nums[left] + nums[right] > nums[i]:
                    res += right - left
                    right -= 1
                else:
                    left += 1
        return res

    def triangleNumber2(self, nums: List[int]) -> int:
        """二分查找O(n^2logn)"""
        n = len(nums)
        nums.sort()
        res = 0
        for i in range(n - 2):
            for j in range(i + 1, n - 1):
                upper = nums[i] + nums[j]
                pos = bisect_left(nums, upper, j + 1)
                res += pos - j - 1
        return res

    def triangleNumber3(self, nums: List[int]) -> int:
        """二分查找O(n^2logn) 要求任意两边长度都不能相等"""
        n = len(nums)
        nums.sort()
        res = 0
        for i in range(n - 2):
            for j in range(i + 1, n - 1):
                if nums[i] == nums[j]:
                    continue
                upper = nums[i] + nums[j]
                pos1 = bisect_right(nums, nums[j], j + 1)
                pos2 = bisect_left(nums, upper, j + 1)
                res += pos2 - pos1
        return res


assert Solution().triangleNumber([2, 2, 3, 4]) == 3
assert Solution().triangleNumber2([2, 2, 3, 4]) == 3
assert Solution().triangleNumber3([2, 2, 3, 4]) == 2
