# 446. 等差数列划分 II - 子序列
# https://leetcode.cn/problems/arithmetic-slices-ii-subsequence/description/
# 如果一个序列中 至少有三个元素 ，并且任意两个相邻元素之差相同，则称该序列为等差序列。
# !求等差数列子序列个数.
# n<=1000
#
# dp[i][diff] 表示以 nums[i] 结尾的公差为 diff 的等差数列个数.

from typing import List
from collections import defaultdict


class Solution:
    def numberOfArithmeticSlices(self, nums: List[int]) -> int:
        """dp[i][diff]表示以nums[i]结尾的公差为diff的等差数列个数."""
        n = len(nums)
        dp = [defaultdict(int) for _ in range(n)]
        res = 0
        for i, v in enumerate(nums):
            for j in range(i):
                diff = v - nums[j]
                count = dp[j][diff]
                dp[i][diff] += count + 1
                res += count
        return res

    def numberOfArithmeticSlices2(self, nums: List[int]) -> int:
        n = len(nums)
        dp = defaultdict(int)
        res = 0
        for i in range(1, n):
            for j in range(i):
                diff = nums[i] - nums[j]
                res += dp[(j, diff)]
                dp[(i, diff)] += dp[(j, diff)] + 1
        return res
