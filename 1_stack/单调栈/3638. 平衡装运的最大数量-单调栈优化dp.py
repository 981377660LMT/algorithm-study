# 3638. 平衡装运的最大数量-单调栈优化dp
# https://leetcode.cn/problems/maximum-balanced-shipments/description/
#
# 给你一个长度为 n 的整数数组 weight，表示按直线排列的 n 个包裹的重量。
# 装运 定义为包裹的一个连续子数组。
# 如果一个装运满足以下条件，则称其为 平衡装运：最后一个包裹的重量 严格小于 该装运中所有包裹中 最大重量 。
# 选择若干个 不重叠 的连续平衡装运，并满足 每个包裹最多出现在一次装运中（部分包裹可以不被装运）。
# 返回 可以形成的平衡装运的最大数量 。
#
# dp[i]: 前 i 个包裹最大平衡装运数量.

from typing import List


def getLeftGreater(nums: List[int]) -> List[int]:
    n = len(nums)
    leftGreater = [-1] * n
    stack = []
    for i in range(n):
        while stack and nums[stack[-1]] <= nums[i]:
            stack.pop()
        leftGreater[i] = stack[-1] if stack else -1
        stack.append(i)
    return leftGreater


class Solution:
    def maxBalancedShipments(self, weight: List[int]) -> int:
        n = len(weight)
        leftGreater = getLeftGreater(weight)
        dp = [0] * (n + 1)
        for i in range(1, n + 1):
            dp[i] = dp[i - 1]
            if leftGreater[i - 1] != -1:
                dp[i] = max(dp[i], dp[leftGreater[i - 1]] + 1)
        return dp[n]
