from typing import List


class Solution:
    def findDisappearedNumbers(self, nums: List[int]) -> List[int]:
        for v in nums:
            pos = abs(v) - 1
            nums[pos] = -abs(nums[pos])
        return [i + 1 for i, v in enumerate(nums) if v > 0]
