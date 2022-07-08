# /**
#  *
#  * @param nums
#  * @description
#  * 给定一个非空整数数组，找到使所有数组元素相等所需的最小移动数，
#  * 其中每次移动可将选定的一个元素加1或减1。
#  * 您可以假设数组的长度最多为10000。
#  */
from typing import List


class Solution:
    def minMoves2(self, nums: List[int]) -> int:
        """462. 最少移动次数使数组元素相等 II-一维曼哈顿距离和最小"""
        nums.sort()
        mid = nums[len(nums) >> 1]
        return sum(abs(num - mid) for num in nums)
