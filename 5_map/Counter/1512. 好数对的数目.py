from typing import List
from collections import Counter
from math import comb


# 如果一组数字 (i,j) 满足 nums[i] == nums[j] 且 i < j ，就可以认为这是一组 好数对 。
class Solution:
    def numIdenticalPairs(self, nums: List[int]) -> int:
        return sum(comb(v, 2) for v in Counter(nums).values())

