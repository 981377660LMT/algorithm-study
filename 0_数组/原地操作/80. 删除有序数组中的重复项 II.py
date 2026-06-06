from typing import List


class Solution:
    def removeDuplicates(self, nums: List[int], k=2) -> int:
        if len(nums) <= k:
            return len(nums)
        slow = k
        for fast in range(k, len(nums)):
            if nums[slow - k] != nums[fast]:
                nums[slow] = nums[fast]
                slow += 1
        return slow
