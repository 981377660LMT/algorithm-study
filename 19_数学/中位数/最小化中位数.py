from typing import List

# 分成等大小的两个数组，使得中位数差最小
class Solution:
    def solve(self, nums: List[int]):
        nums.sort()
        n = len(nums)
        return nums[n // 2] - nums[n // 2 - 1]


print(Solution().solve(nums=[1, 9, 7, 4, 3, 6]))
