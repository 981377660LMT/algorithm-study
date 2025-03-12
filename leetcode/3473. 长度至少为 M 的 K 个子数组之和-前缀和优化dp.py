# 3473. 长度至少为 M 的 K 个子数组之和-前缀和优化dp
# https://leetcode.cn/problems/sum-of-k-subarrays-with-length-at-least-m/description/
# 给你一个整数数组 nums 和两个整数 k 和 m。
# 返回数组 nums 中 k 个不重叠子数组的 最大 和，其中每个子数组的长度 至少 为 m。
#
# dp[i][rest]表示nums[i:]选rest个长度至少为M的不重叠子数组的最大和
#

from itertools import accumulate
from typing import List

INF = int(1e18)


class Solution:
    def maxSum(self, nums: List[int], k: int, m: int) -> int:
        n = len(nums)
        presum = list(accumulate(nums, initial=0))
        dp = [0] * (n + 1)
        for i in range(1, k + 1):
            ndp = [-INF] * (n + 1)
            max_ = -INF
            for j in range(i * m, n - (k - i) * m + 1):
                max_ = max(max_, dp[j - m] - presum[j - m])
                ndp[j] = max(ndp[j - 1], max_ + presum[j])
            dp = ndp
        return dp[n]
