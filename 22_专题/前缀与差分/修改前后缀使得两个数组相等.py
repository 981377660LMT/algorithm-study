# https://www.codechef.com/START30A/problems-old/PRESUFOP
# Prefix Suffix Operations
# 前后缀+1使得两个数组相等

# 贪心

from typing import List

# 设下标 i 加了 xi 次 prefix，加了 yi 次 suffix，
# 那答案就是 x1 + yn，要求 xi 不增，yi 不减。

# 然后我们把加 1 ~ n 也看作加 prefix，
# 这样有 x1 = diff1，y1 = 0。
# 所以我们要最小化 yn，
# 那么 x2 = min(x1, diff2), y2 = diff2 - x2，
# x3 = min(x2, diff3)，y3 = diff3 - x3 这样贪心。


def prefixSuffixOperations(nums1: List[int], nums2: List[int]) -> int:
    """
    可以给nums1的前缀或者后缀全部+1 求将两个数组变为相等的最小操作次数
    如果不可能返回-1
    """
    n = len(nums1)
    diff = [b - a for a, b in zip(nums1, nums2)]
    if any(d < 0 for d in diff):
        return -1

    preAdd, sufAdd = [diff[0]] * n, [0] * n  # 每个数前缀加的次数，后缀加的次数
    for i in range(1, n):
        preAdd[i] = min(preAdd[i - 1], diff[i] - sufAdd[i - 1])
        sufAdd[i] = diff[i] - preAdd[i]
        if preAdd[i] < 0 or sufAdd[i] < 0:
            return -1

    return max(preAdd) + max(sufAdd)


assert prefixSuffixOperations([2, 3, 5, 1, 2], [4, 3, 6, 2, 3]) == 3
assert prefixSuffixOperations([0, 0, 0, 0], [1, 2, 2, 1]) == 2
assert prefixSuffixOperations([1, 2, 3], [1, 2, 2]) == -1
