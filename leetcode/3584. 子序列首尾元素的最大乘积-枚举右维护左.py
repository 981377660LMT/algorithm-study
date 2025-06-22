# 3584. 子序列首尾元素的最大乘积
# https://leetcode.cn/problems/maximum-product-of-first-and-last-elements-of-a-subsequence/
# 给你一个整数数组 nums 和一个整数 m。
# !返回任意大小为 m 的 子序列 中首尾元素乘积的最大值。
# 子序列 是可以通过删除原数组中的一些元素（或不删除任何元素），且不改变剩余元素顺序而得到的数组。
#
# !nums 的任意下标相差至少为 m−1 的两数之积的最大值

from typing import List


INF = int(1e20)


class Solution:
    def maximumProduct(self, nums: List[int], m: int) -> int:
        res = -INF
        preMin, preMax = INF, -INF
        for i in range(m - 1, len(nums)):
            pre = nums[i - (m - 1)]
            preMin = min(preMin, pre)
            preMax = max(preMax, pre)
            cur = nums[i]
            res = max(res, cur * preMin, cur * preMax)
        return res
