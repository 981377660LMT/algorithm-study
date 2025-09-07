# 3676. 碗子数组的数目
# https://leetcode.cn/problems/count-bowl-subarrays/solutions/3774503/dan-diao-zhan-by-tsreaper-oogs/
# 给你一个整数数组 nums，包含 互不相同 的元素。
#
# nums 的一个子数组 nums[l...r] 被称为 碗（bowl），如果它满足以下条件：
#
# 子数组的长度至少为 3。也就是说，r - l + 1 >= 3。
# 其两端元素的 最小值 严格大于 中间所有元素的 最大值。也就是说，min(nums[l], nums[r]) > max(nums[l + 1], ..., nums[r - 1])。
# 返回 nums 中 碗 子数组的数量。
#
# 子数组 是数组中连续的元素序列。
from typing import List

from 每个元素作为最值的影响范围 import getRange


class Solution:
    def bowlSubarrays(self, nums: List[int]) -> int:
        n = len(nums)
        ranges = getRange(nums, isMax=True, isLeftStrict=True, isRightStrict=False)
        lefts = [-1] * n
        rights = [-1] * n
        for i, (l, r) in enumerate(ranges):
            if l - 1 >= 0:
                lefts[i] = l - 1
            if r + 1 < n:
                rights[i] = r + 1

        res = 0
        for i in range(n):
            if lefts[i] != -1 and i - lefts[i] >= 2:
                res += 1
            if rights[i] != -1 and rights[i] - i >= 2:
                res += 1
        return res
