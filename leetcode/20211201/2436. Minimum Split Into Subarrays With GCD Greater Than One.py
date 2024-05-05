from math import gcd
from typing import List


# gcd的单调性
class Solution:
    def minimumSplits(self, nums: List[int]) -> int:
        """分割数组成最少的子数组,使得每个子数组的最大公约数大于1"""
        n, res, preGcd = len(nums), 1, nums[0]
        for i in range(1, n):
            curGcd = gcd(preGcd, nums[i])
            if curGcd == 1:
                res += 1
                curGcd = nums[i]
            preGcd = curGcd
        return res


assert Solution().minimumSplits([12, 6, 3, 14, 8]) == 2
