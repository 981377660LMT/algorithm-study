from typing import List


class Solution:
    def check(self, nums: List[int]) -> bool:
        n = len(nums)
        down = 0
        for i in range(n):
            if nums[i] > nums[(i + 1) % n]:
                down += 1
            if down > 1:
                return False
        return True
