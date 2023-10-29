# https://leetcode.cn/problems/longest-increasing-subsequence/
# O(n^2) dp 解法


from typing import List


class Solution:
    def lengthOfLIS(self, nums: List[int]) -> int:
        n = len(nums)
        dp = [1] * n  # dp[i] 表示以 nums[i] 结尾的最长上升子序列的长度
        for i in range(1, n):
            for j in range(i):
                if nums[j] < nums[i]:
                    dp[i] = max(dp[j] + 1, dp[i])
        return max(dp)
