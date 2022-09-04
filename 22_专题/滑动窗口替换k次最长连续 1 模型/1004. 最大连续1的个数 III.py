# 1004. 最大连续1的个数 III
from typing import List, Sequence, TypeVar

T = TypeVar("T")


def fix(nums: Sequence[T], need: T, k: int) -> int:
    """改变最多 k 个 字符,求 nums 中最大连续 need 的个数"""
    n, left, res = len(nums), 0, 0
    for right in range(n):
        if nums[right] != need:
            k -= 1
        while k < 0:
            if nums[left] != need:
                k += 1
            left += 1
        res = max(res, right - left + 1)
    return res


class Solution:
    def longestOnes(self, nums: List[int], k: int) -> int:
        return fix(nums, 1, k)
