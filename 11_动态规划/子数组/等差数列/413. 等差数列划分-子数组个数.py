# 413. 等差数列划分
# https://leetcode.cn/problems/arithmetic-slices/
# 如果一个数列 至少有三个元素 ，并且任意两个相邻元素之差相同，则称该数列为等差数列。
# !求等差数列子数组个数.
# n<=1e5


from typing import List


class Solution:
    def numberOfArithmeticSlices(self, nums: List[int]) -> int:
        n = len(nums)
        dp = [0] * (n + 1)
        for i in range(2, n):
            if nums[i] - nums[i - 1] == nums[i - 1] - nums[i - 2]:
                dp[i] = dp[i - 1] + 1
        return sum(dp)
