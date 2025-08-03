# 100743. 使数组平衡的最少移除数目
#
# 给你一个整数数组 nums 和一个整数 k。
# 如果一个数组的 最大 元素的值 至多 是其 最小 元素的 k 倍，则该数组被称为是 平衡 的。
# 你可以从 nums 中移除 任意 数量的元素，但不能使其变为 空 数组。
# 返回为了使剩余数组平衡，需要移除的元素的 最小 数量。
#
# !枚举最大值 r，用双指针找到最小值 l 最小可以是多少

from typing import List


class Solution:
    def minRemoval(self, nums: List[int], k: int) -> int:
        nums = sorted(nums)
        maxKeep = 0
        left = 0
        for right, v in enumerate(nums):
            while nums[left] * k < v:
                left += 1
            maxKeep = max(maxKeep, right - left + 1)
        return len(nums) - maxKeep
