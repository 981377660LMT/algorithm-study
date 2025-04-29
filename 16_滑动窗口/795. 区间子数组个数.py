# 795. 区间子数组个数
# https://leetcode.cn/problems/number-of-subarrays-with-bounded-maximum/description/
# 类似 https://leetcode.cn/problems/count-subarrays-with-fixed-bounds/
#
# 给你一个整数数组 nums 和两个整数：left 及 right 。
# 找出 nums 中连续、非空且其中最大元素在范围 [left, right] 内的子数组，并返回满足条件的子数组的个数。

from typing import List


class Solution:
    def numSubarrayBoundedMax(self, nums: List[int], left: int, right: int) -> int:
        res = 0
        pos1, pos2 = -1, -1
        for i, v in enumerate(nums):
            if v > right:
                pos1 = i
            if v >= left:
                pos2 = i
            res += pos2 - pos1
        return res
