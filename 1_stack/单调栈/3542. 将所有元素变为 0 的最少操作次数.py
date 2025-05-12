# 3542. 将所有元素变为 0 的最少操作次数
# https://leetcode.cn/problems/minimum-operations-to-convert-all-elements-to-zero/description/
# 给你一个大小为 n 的 非负 整数数组 nums 。你的任务是对该数组执行若干次（可能为 0 次）操作，使得 所有 元素都变为 0。
# 在一次操作中，你可以选择一个子数组 [i, j]（其中 0 <= i <= j < n），将该子数组中所有 最小的非负整数 的设为 0。
# 返回使整个数组变为 0 所需的最少操作次数。


from typing import List

from 每个元素作为最值的影响范围 import getRange


class Solution:
    def minOperations(self, nums: List[int]) -> int:
        ranges = getRange(nums, isMax=False, isLeftStrict=False, isRightStrict=False)
        s = set()
        for v, lr in zip(nums, ranges):
            if v > 0:
                s.add(lr)
        return len(s)
