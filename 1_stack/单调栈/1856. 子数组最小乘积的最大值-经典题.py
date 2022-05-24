from itertools import accumulate
from typing import List
from 每个元素作为最值的影响范围 import getRange

MOD = int(1e9 + 7)


# 一个数组的 最小乘积 定义为这个数组中 最小值 乘以 数组的 和 。

# 思路：
# 我们可以考虑枚举数组中每个值n，并且以n作为子数组中的最小值，再乘以这个子数组的和
# 给定n，如何找到以n为最小值的子数组边界？
# 对于数组中的每一个元素，分别找到左边、右边第一个比它小的元素的位置


class Solution:
    def maxSumMinProduct(self, nums: List[int]) -> int:
        res = 0
        preSum = [0] + list(accumulate(nums))
        minRange = getRange(nums, isLeftStrict=False, isRightStrict=False, isMax=False)
        for i, num in enumerate(nums):
            left, right = minRange[i]
            res = max((preSum[right + 1] - preSum[left]) * num, res)
        return res % MOD


print(Solution().maxSumMinProduct(nums=[1, 2, 3, 2]))
print(Solution().maxSumMinProduct(nums=[2, 3, 3, 1, 2]))
print(Solution().maxSumMinProduct(nums=[3, 1, 5, 6, 4, 2]))
# 输出：14
# 解释：最小乘积的最大值由子数组 [2,3,2] （最小值是 2）得到。
# 2 * (2+3+2) = 2 * 7 = 14 。
