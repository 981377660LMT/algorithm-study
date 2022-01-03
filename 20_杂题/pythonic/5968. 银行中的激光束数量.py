from typing import List


class Solution:
    def numberOfBeams(self, bank: List[str]) -> int:
        nums = [c for row in bank if (c := row.count('1')) != 0]
        return sum(a * b for a, b in zip(nums, nums[1:]))
