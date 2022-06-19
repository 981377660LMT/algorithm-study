# 1004. 最大连续1的个数 III
from typing import List, Sequence, TypeVar

T = TypeVar('T')


def fix(raw: Sequence[T], need: T, k: int) -> int:
    """改变最多 k 个 字符,求 raw 中最大连续 need 的个数

    Args:
        raw (T): 源字符串
        need (T): 关心的字符
        k (int): 可替换k次

    Returns:
        int: 最大连续长度
    """
    left, res = 0, 0
    for right in range(len(raw)):
        if raw[right] != need:
            k -= 1
        while k < 0:
            if raw[left] != need:
                k += 1
            left += 1
        res = max(res, right - left + 1)
    return res


class Solution:
    def longestOnes(self, nums: List[int], k: int) -> int:
        return fix(nums, 1, k)

