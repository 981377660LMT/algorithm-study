# 3471. 找出最大的几近缺失整数
# https://leetcode.cn/problems/find-the-largest-almost-missing-integer/description/
# 给你一个整数数组 nums 和一个整数 k 。
# 如果整数 x 恰好仅出现在 nums 中的一个大小为 k 的子数组中，则认为 x 是 nums 中的几近缺失（almost missing）整数。
# 返回 nums 中 最大的几近缺失 整数，如果不存在这样的整数，返回 -1 。
# 子数组 是数组中的一个连续元素序列。


from typing import List
from collections import Counter


class Solution:
    def largestInteger(self, nums: List[int], k: int) -> int:
        if k == len(nums):
            return max(nums)

        if k == 1:
            res = -1
            counter = Counter(nums)
            for k, v in counter.items():
                if v == 1:
                    res = max(res, k)
            return res

        def f(arr: List[int], v: int):
            return -1 if v in arr else v

        return max(f(nums[1:], nums[0]), f(nums[:-1], nums[-1]))
