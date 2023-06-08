# 求邻位交换的最小次数.
# 使得数组的第一个数字等于 1 且最后一个数字等于 n.
from typing import List


class Solution:
    def semiOrderedPermutation(self, nums: List[int]) -> int:
        n = len(nums)
        i1 = nums.index(1)
        i2 = next((i for i in range(n - 1, -1, -1) if nums[i] == n))
        return i1 + n - 1 - i2 - (i1 > i2)


assert Solution().semiOrderedPermutation(nums=[2, 4, 1, 3]) == 3
