# 3785. 避免禁用值的最小交换次数
# https://leetcode.cn/problems/minimum-swaps-to-avoid-forbidden-values/description/
# 你可以执行以下操作任意次（包括零次）：
# 选择两个 不同 下标 i 和 j，然后交换 nums[i] 和 nums[j]。
# 返回使得对于每个下标 i，nums[i] 不等于 forbidden[i] 所需的 最小 交换次数。如果无论如何都无法满足条件，返回 -1。
#
# !经典问题：给一个序列，每次可以删掉一个元素，或选择两个不同的元素同时删掉，删多少次才能让序列变为空

from typing import List
from collections import Counter


class Solution:
    def minSwaps(self, nums: List[int], forbidden: List[int]) -> int:
        n = len(nums)
        total = Counter(nums) + Counter(forbidden)
        if any(c > n for c in total.values()):
            return -1
        same = Counter(x for x, y in zip(nums, forbidden) if x == y)
        sum_ = same.total()
        max_ = max(same.values(), default=0)
        return max(max_, (sum_ + 1) // 2)
