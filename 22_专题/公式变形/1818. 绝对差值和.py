"""
数组 nums1 和 nums2 的 绝对差值和 定义为所有 |nums1[i] - nums2[i]|（0 <= i < n）的 总和（下标从 0 开始）。
你可以选用 nums1 中的 任意一个 元素来替换 nums1 中的 至多 一个元素，以 最小化 绝对差值和。
在替换数组 nums1 中最多一个元素 之后 ，返回最小绝对差值和。
因为答案可能很大，所以需要对 1e9 + 7 取余 后返回。
"""

# 警惕sort的可变性

from bisect import bisect_right
from typing import List

MOD = int(1e9 + 7)
INF = int(1e18)


class Solution:
    def minAbsoluteSumDiff(self, nums1: List[int], nums2: List[int]) -> int:
        n = len(nums1)
        sl = sorted(nums1)
        res = sum(abs(nums1[i] - nums2[i]) for i in range(n))
        max_ = -INF
        for i in range(n):
            pos = bisect_right(sl, nums2[i])
            if pos < n:
                max_ = max(max_, abs(nums1[i] - nums2[i]) - abs(sl[pos] - nums2[i]))
            if pos - 1 >= 0:
                max_ = max(max_, abs(nums1[i] - nums2[i]) - abs(sl[pos - 1] - nums2[i]))
        return (res - max_) % MOD
