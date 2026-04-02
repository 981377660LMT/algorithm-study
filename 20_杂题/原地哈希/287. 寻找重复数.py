from typing import List


class Solution:
    def findDuplicate(self, nums: List[int]) -> int:
        for v in nums:
            pos = abs(v) - 1
            if nums[pos] < 0:
                return pos + 1
            nums[pos] = -abs(nums[pos])
        return -1
