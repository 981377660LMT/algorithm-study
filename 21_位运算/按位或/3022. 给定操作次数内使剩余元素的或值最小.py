# 3022. 给定操作次数内使剩余元素的或值最小
# https://leetcode.cn/problems/minimize-or-of-remaining-elements-using-operations/description/
# 给你一个下标从 0 开始的整数数组 nums 和一个整数 k 。
# 一次操作中，你可以选择 nums 中满足 0 <= i < nums.length - 1 的一个下标 i ，
# 并将 nums[i] 和 nums[i + 1] 替换为数字 nums[i] & nums[i + 1] ，其中 & 表示按位 AND 操作。
# !请你返回 至多 k 次操作以内，使 nums 中所有剩余元素按位 OR 结果的 最小值 。

from typing import List


class Solution:
    def minOrAfterOperations(self, nums: List[int], k: int) -> int:
        ...
