from itertools import accumulate
from typing import List


MOD = int(1e9 + 7)


# 一个数组的 最小乘积 定义为这个数组中 最小值 乘以 数组的 和 。

# 思路：
# 我们可以考虑枚举数组中每个值n，并且以n作为子数组中的最小值，再乘以这个子数组的和
# 给定n，如何找到以n为最小值的子数组边界？
# 对于数组中的每一个元素，分别找到左边、右边第一个比它小的元素的位置


class Solution:
    def maxSumMinProduct(self, nums: List[int]) -> int:
        """看每个数作为非严格最小值的影响范围"""
        n = len(nums)
        preSum = [0] + list(accumulate(nums))

        # 右边第一个比它小的元素下标
        rightSmaller = [n] * n
        stack = []
        for i in range(len(nums)):
            while stack and nums[i] < nums[stack[-1]]:
                rightSmaller[stack.pop()] = i
            stack.append(i)

        leftSmaller = [-1] * n
        stack = []
        for i in range(len(nums) - 1, -1, -1):
            while stack and nums[i] < nums[stack[-1]]:
                leftSmaller[stack.pop()] = i
            stack.append(i)

        res = 0
        for i in range(len(nums)):
            left = leftSmaller[i]
            right = rightSmaller[i]
            res = max(res, nums[i] * (preSum[right] - preSum[left + 1]))
        return res % MOD


print(Solution().maxSumMinProduct(nums=[1, 2, 3, 2]))
print(Solution().maxSumMinProduct(nums=[2, 3, 3, 1, 2]))
print(Solution().maxSumMinProduct(nums=[3, 1, 5, 6, 4, 2]))
# 输出：14
# 解释：最小乘积的最大值由子数组 [2,3,2] （最小值是 2）得到。
# 2 * (2+3+2) = 2 * 7 = 14 。
