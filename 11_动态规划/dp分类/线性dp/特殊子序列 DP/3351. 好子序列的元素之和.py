# 3351. 好子序列的元素之和
# https://leetcode.cn/problems/sum-of-good-subsequences/description/
# !给你一个整数数组 nums。好子序列 的定义是：子序列中任意 两个 连续元素的绝对差 恰好 为 1。
# 子序列 是指可以通过删除某个数组的部分元素（或不删除）得到的数组，并且不改变剩余元素的顺序。
# 返回 nums 中所有 可能存在的 好子序列的 元素之和。
# 因为答案可能非常大，返回结果需要对 109 + 7 取余。
# 注意，长度为 1 的子序列默认为好子序列。

from collections import defaultdict
from typing import List


MOD = int(1e9 + 7)


class Solution:
    def sumOfGoodSubsequences(self, nums: List[int]) -> int:
        dpCount, dpSum = defaultdict(int), defaultdict(int)
        for v in nums:
            c = dpCount[v - 1] + dpCount[v + 1] + 1  # 单独、在v-1之后、在v+1之后
            dpSum[v] = (dpSum[v] + dpSum[v - 1] + dpSum[v + 1] + c * v) % MOD
            dpCount[v] = (dpCount[v] + c) % MOD
        return sum(dpSum.values()) % MOD
