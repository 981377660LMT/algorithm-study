from math import gcd
from typing import List

# 你可以从 nums 中删除任意数目的元素。
# 请你返回使 nums 中 最小 元素可以整除 numsDivide 中所有元素的 最少 删除次数。
# 如果无法得到这样的元素，返回 -1 。


class Solution:
    def minOperations(self, nums: List[int], numsDivide: List[int]) -> int:
        """求出最小的被gcd整除的数 删除所有小于这个数的数"""
        gcd_ = gcd(*numsDivide)
        min_ = min((num for num in nums if gcd_ % num == 0), default=-1)
        if min_ == -1:
            return -1
        return sum(num < min_ for num in nums)


print(Solution().minOperations(nums=[2, 3, 2, 4, 3], numsDivide=[9, 6, 9, 3, 15]))
