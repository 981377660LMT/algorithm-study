from heapq import nlargest
from typing import List


# 2 <= nums.length <= 500
# 1 <= nums[i] <= 10^3
class Solution:
    def maxProduct(self, nums: List[int]) -> float:
        a, b = nlargest(2, nums)
        return (a - 1) * (b - 1)

    def maxProduct2(self, nums: List[int]) -> float:
        a, b = 0, 0
        for num in nums:
            if num > a:
                a, b = num, a
            elif num > b:
                b = num
        return (a - 1) * (b - 1)


print(Solution().maxProduct([3, 4, 5, 2]))
