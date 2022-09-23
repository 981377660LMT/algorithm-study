# 1004. 最大连续1的个数 III
# fix模型/fixK模型

from typing import List, Sequence, TypeVar

T = TypeVar("T")


def fix(seq: Sequence[T], target: T, k: int) -> int:
    """改变最多 k 个 字符,求 seq 中最大连续 target 的个数"""
    n, left, res = len(seq), 0, 0
    for right in range(n):
        if seq[right] != target:
            k -= 1
        while k < 0:
            if seq[left] != target:
                k += 1
            left += 1
        res = max(res, right - left + 1)
    return res


class Solution:
    def longestOnes(self, nums: List[int], k: int) -> int:
        return fix(nums, 1, k)
