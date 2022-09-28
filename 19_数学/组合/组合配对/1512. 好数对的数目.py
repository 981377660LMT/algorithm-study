from typing import List
from collections import Counter
from math import comb


# !如果一组数字 (i,j) 满足 nums[i] == nums[j] 且 i < j ，就可以认为这是一组 好数对 。
# O(n)
class Solution:
    def numIdenticalPairs(self, nums: List[int]) -> int:
        return sum(comb(v, 2) for v in Counter(nums).values())

    def numIdenticalPairs2(self, nums: List[int]) -> int:
        """一遍遍历"""
        counter = dict()
        res = 0
        for num in nums:
            res += counter.get(num, 0)
            counter[num] = counter.get(num, 0) + 1
        return res
