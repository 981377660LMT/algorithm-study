from typing import List

# Given a list of integers nums,
# return the length of the shortest sublist in nums which if sorted would make nums sorted in ascending order.
class Solution:
    def solve(self, nums: List[int]) -> int:
        if not nums:
            return 0
        """从左往右找最后的下降，从右往左找最后的上升"""
        n = len(nums)

        leftMax = nums[0]
        end = -1
        for i in range(1, n):
            if nums[i] < leftMax:
                end = i
            leftMax = max(leftMax, nums[i])
        if end == -1:
            return 0

        rightMin = nums[-1]
        start = -1
        for i in range(n - 2, -1, -1):
            if nums[i] > rightMin:
                start = i
            rightMin = min(rightMin, nums[i])
        if start == -1:
            return 0

        return end - start + 1

