# 3599. 划分数组得到最小 XOR
# https://leetcode.cn/problems/partition-array-to-minimize-xor/description/
#
# 给你一个整数数组 nums 和一个整数 k。
# 你的任务是将 nums 分成 k 个非空的 子数组 。对每个子数组，计算其所有元素的按位 XOR 值。
# 返回这 k 个子数组中 最大 XOR 的 最小值 。
# 子数组 是数组中连续的 非空 元素序列。
#
# !dp[i][j] 表示前i个数分成j个子数组的最大xor的最小值。
# n<=250

from typing import List


INF = int(1e20)


def min2(a: int, b: int) -> int:
    return a if a < b else b


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def minXor(self, nums: List[int], k: int) -> int:
        n = len(nums)
        prexor = [0] * (n + 1)
        for i in range(n):
            prexor[i + 1] = prexor[i] ^ nums[i]

        dp = [INF] * (n + 1)
        dp[0] = 0
        for _ in range(k):
            ndp = [INF] * (n + 1)
            for i in range(1, n + 1):
                for j in range(i):
                    ndp[i] = min2(ndp[i], max2(dp[j], prexor[i] ^ prexor[j]))
            dp = ndp

        return dp[n]
