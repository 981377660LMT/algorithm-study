# Return the sum of every concatenation of every pair of numbers in nums.
from itertools import product
from typing import Counter


class Solution:
    def solve(self, nums):
        const = len(nums) + sum(10 ** len(str(x)) for x in nums)
        return sum(const * num for num in nums)


print(Solution().solve(nums=[10, 2]))

# We have the following concatenations:

# nums[0] + nums[0] = 1010
# nums[0] + nums[1] = 102
# nums[1] + nums[0] = 210
# nums[1] + nums[1] = 22
# And its sum is 1344
