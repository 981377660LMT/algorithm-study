# 数组是否值域连续
# 2229. Check if an Array Is Consecutive
from typing import List


class Solution:
    def isConsecutive(self, nums: List[int]) -> bool:
        return max(nums) - min(nums) == len(nums) - 1 and len(set(nums)) == len(nums)

