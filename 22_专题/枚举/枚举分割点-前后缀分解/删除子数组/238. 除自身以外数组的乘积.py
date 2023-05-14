# 238. 除自身以外数组的乘积


from typing import List
from productWithoutOne import productWithoutOne


class Solution:
    def productExceptSelf(self, nums: List[int]) -> List[int]:
        return productWithoutOne(nums, lambda: 1, lambda x, y: x * y)


assert Solution().productExceptSelf(nums=[1, 2, 3, 4]) == [24, 12, 8, 6]
