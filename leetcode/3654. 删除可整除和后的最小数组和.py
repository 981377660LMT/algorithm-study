# 3654. 删除可整除和后的最小数组和
# https://leetcode.cn/problems/minimum-sum-after-divisible-sum-deletions/description/
#
# 给你一个整数数组 nums 和一个整数 k。
# 你可以 多次 选择 连续 子数组 nums，其元素和可以被 k 整除，并将其删除；
# 每次删除后，剩余元素会填补空缺。
# 返回在执行任意次数此类删除操作后，nums 的最小可能 和。
#

from typing import List

INF = int(1e18)


class Solution:
    def minArraySum(self, nums: List[int], k: int) -> int:
        n = len(nums)
        curSum = 0
        pre = {0: 0}
        dp = [0] * (n + 1)
        for i, v in enumerate(nums, 1):
            curSum = (curSum + v) % k
            res1 = dp[i - 1] + v  # 不删除当前元素
            res2 = dp[pre[curSum]] if curSum in pre else INF  # 删除当前元素
            dp[i] = min(res1, res2)
            pre[curSum] = i
        return dp[n]
