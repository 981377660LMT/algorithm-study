# /**
#  * @param {number[]} nums
#  * @return {number}
#  * @description 有点像 53. 最大子序和
#  * 我们只要记录前i的最小值, 和最大值
#  */
from typing import List


class Solution:
    def maxProduct(self, nums: List[int]) -> int:
        """子数组：前面一截取还是全不取"""
        res = min_ = max_ = nums[0]
        for num in nums[1:]:
            min_, max_ = min(min_ * num, max_ * num, num), max(min_ * num, max_ * num, num)
            res = max(res, max_)
        return res

