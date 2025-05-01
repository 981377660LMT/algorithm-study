# 2248. 多个数组求交集

from typing import List


class Solution:
    def intersection(self, nums: List[List[int]]) -> List[int]:
        if not nums:
            return []
        common = set(nums[0])
        for arr in nums[1:]:
            common.intersection_update(set(arr))
        return sorted(common)
