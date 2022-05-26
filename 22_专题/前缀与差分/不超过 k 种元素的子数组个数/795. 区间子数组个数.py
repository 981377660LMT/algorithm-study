# 给定一个元素都是正整数的数组 A ，正整数 L  以及  R (L <= R)。
# 求连续、非空且其中最大元素满足大于等于 L  小于等于 R 的子数组个数。
from typing import List


def atMaxK(nums: List[int], k: int) -> int:
    """元素最大值个数<=k的子数组个数 dp求出"""
    res, dp = 0, 0
    for num in nums:
        dp = (dp + 1) if num <= k else 0
        res += dp
    return res


class Solution:
    def numSubarrayBoundedMax(self, nums: List[int], left: int, right: int) -> int:
        return atMaxK(nums, right) - atMaxK(nums, left - 1)

