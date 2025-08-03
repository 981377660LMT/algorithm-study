# 3637. 三段式数组 I
# https://leetcode.cn/problems/trionic-array-i/description/
# 给你一个长度为 n 的整数数组 nums。
#
# 如果存在索引 0 < p < q < n − 1，使得数组满足以下条件，则称其为 三段式数组（trionic）：
#
# nums[0...p] 严格 递增，
# nums[p...q] 严格 递减，
# nums[q...n − 1] 严格 递增。
# 如果 nums 是三段式数组，返回 true；否则，返回 false。

from typing import List


class Solution:
    def isTrionic(self, nums: List[int]) -> bool:
        n = len(nums)

        s0 = 1
        while s0 < n and nums[s0 - 1] < nums[s0]:
            s0 += 1
        if s0 == 1:
            return False

        s1 = s0
        while s1 < n and nums[s1 - 1] > nums[s1]:
            s1 += 1
        if s1 == s0 or s1 == n:
            return False

        s2 = s1
        while s2 < n and nums[s2 - 1] < nums[s2]:
            s2 += 1
        return s2 == n
