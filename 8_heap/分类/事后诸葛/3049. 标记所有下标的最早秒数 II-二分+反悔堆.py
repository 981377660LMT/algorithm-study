# 3049. 标记所有下标的最早秒数 II-二分+反悔堆
# https://leetcode.cn/problems/earliest-second-to-mark-indices-ii/description/
#
# 给你两个下标从 1 开始的整数数组 nums 和 changeIndices ，数组的长度分别为 n 和 m 。
# 一开始，nums 中所有下标都是未标记的，你的任务是标记 nums 中 所有 下标。
# 从第 1 秒到第 m 秒（包括 第 m 秒），对于每一秒 s ，你可以执行以下操作 之一 ：
# - 选择范围 [1, n] 中的一个下标 i ，并且将 nums[i] 减少 1 。
# - 将 nums[changeIndices[s]] 设置成任意的 非负 整数。
# - 选择范围 [1, n] 中的一个下标 i ， 满足 nums[i] 等于 0, 并 标记 下标 i 。
# - 什么也不做。
# 请你返回范围 [1, m] 中的一个整数，表示最优操作下，标记 nums 中 所有 下标的 最早秒数 ，如果无法标记所有下标，返回 -1 。
# !TODO
# 题意有点抽象，形象地解释一下：
# 你有 n门课程需要考试，第 i门课程需要用 nums[i] 天复习。同一天只能复习一门课程（慢速复习）。
# 在第 i天，你可以快速搞定第 changeIndices[i] 门课程的复习。
# 你可以在任意一天完成一门课程的考试（前提是复习完成）。考试这一天不能复习。
# 搞定所有课程的复习+考试，至少要多少天？
# !如何权衡哪些课程快速复习，哪些课程慢速复习呢？

from typing import List


class Solution:
    def earliestSecondToMarkIndices(self, nums: List[int], changeIndices: List[int]) -> int:
        ...
