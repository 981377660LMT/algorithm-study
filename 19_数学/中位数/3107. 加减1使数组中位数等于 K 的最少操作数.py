# 3107. 加减1使数组中位数等于 K 的最少操作数
# https://leetcode.cn/problems/minimum-operations-to-make-median-of-array-equal-to-k/description/
# 给你一个整数数组 nums 和一个 非负 整数 k 。一次操作中，你可以选择任一元素 加 1 或者减 1 。
# 请你返回将 nums 中位数 变为 k 所需要的 最少 操作次数。
# 一个数组的中位数指的是数组按非递减顺序排序后最中间的元素。
# 如果数组长度为偶数，我们选择中间两个数的较大值为中位数。
#
# !nums排序后, 需要把中位数左边的数都变成 ≤k 的，右边的数都变成 ≥k 的

from typing import List
from bisect import bisect_right
from itertools import accumulate


class Solution:
    def minOperationsToMakeMedianK(self, nums: List[int], k: int) -> int:
        n = len(nums)
        nums = sorted(nums)
        preSum = [0] + list(accumulate(nums))
        mid = n >> 1

        def solve(to: int) -> int:
            """把[0,mid]变成<=to, [mid+1:]变成>=to"""
            pos = bisect_right(nums, to)
            if pos <= mid:
                return (preSum[mid + 1] - preSum[pos]) - to * (mid + 1 - pos)
            else:
                return to * (pos - mid) - (preSum[pos] - preSum[mid])

        return solve(k)
