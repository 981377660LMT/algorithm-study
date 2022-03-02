from typing import List


MOD = int(1e9 + 7)


# 一个数组的 最小乘积 定义为这个数组中 最小值 乘以 数组的 和 。

# 思路：
# 我们可以考虑枚举数组中每个值n，并且以n作为子数组中的最小值，再乘以这个子数组的和
# 给定n，如何找到以n为最小值的子数组边界？
# 对于数组中的每一个元素，分别找到左边、右边第一个比它小的元素的位置


class Solution:
    def maxSumMinProduct(self, nums: List[int]) -> int:
        """看每个数作为最小值的影响范围"""
        nums = [0] + nums + [0]
        preSum = [0]
        for num in nums:
            preSum.append(preSum[-1] + num)

        # 右边第一个比它小的元素下标
        rightFirstSmaller = [-1] * len(nums)
        stack = []
        for i in range(len(nums)):
            while stack and nums[i] < nums[stack[-1]]:
                rightFirstSmaller[stack.pop()] = i
            stack.append(i)

        leftFirstSmaller = [-1] * len(nums)
        stack = []
        for i in range(len(nums) - 1, -1, -1):
            while stack and nums[i] < nums[stack[-1]]:
                leftFirstSmaller[stack.pop()] = i
            stack.append(i)

        # print(rightFirstSmaller, leftFirstSmaller)
        res = 0
        for i in range(1, len(nums) - 1):
            left = leftFirstSmaller[i]
            right = rightFirstSmaller[i]
            res = max(res, nums[i] * (preSum[right] - preSum[left + 1]))
        return res % MOD


print(Solution().maxSumMinProduct(nums=[1, 2, 3, 2]))
# 输出：14
# 解释：最小乘积的最大值由子数组 [2,3,2] （最小值是 2）得到。
# 2 * (2+3+2) = 2 * 7 = 14 。
