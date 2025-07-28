# 3229. 使数组等于目标数组所需的最少操作次数
# https://leetcode.cn/problems/minimum-operations-to-make-array-equal-to-target/
# 给你两个长度相同的正整数数组 nums 和 target。
# 在一次操作中，你可以选择 nums 的任何子数组，并将该子数组内的每个元素的值增加或减少 1。
# 返回使 nums 数组变为 target 数组所需的 最少 操作次数。
# 等价差分数组于一个数+1，另一个-1
#
# !22_专题/前缀与差分/差分数组/区间操作/Q3. 灯光调整.py

from itertools import pairwise
from typing import List


class Solution:
    def minimumOperations(self, nums: List[int], target: List[int]) -> int:
        diff = [0] + [b - a for a, b in zip(nums, target)]
        diff = [b - a for a, b in pairwise(diff)]
        posSum, negSum = 0, 0  # !如果前面增大/减小了k，后面就可以减小/增大k
        for d in diff:
            if d > 0:
                posSum += d
            elif d < 0:
                negSum += -d
        return max(posSum, negSum)
