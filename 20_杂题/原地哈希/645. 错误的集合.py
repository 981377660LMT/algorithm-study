from typing import List


class Solution:
    def findErrorNums(self, nums: List[int]) -> List[int]:
        dup = -1
        for v in nums:
            pos = abs(v) - 1
            if nums[pos] < 0:
                dup = pos + 1
            else:
                nums[pos] = -abs(nums[pos])
        miss = next(i + 1 for i, v in enumerate(nums) if v > 0)
        return [dup, miss]
