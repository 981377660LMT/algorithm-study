from typing import List

# 请你返回在新数组中下标为 left 到 right （下标从 1 开始）的所有数字和（包括左右端点）。由于答案可能很大，请你将它对 10^9 + 7 取模后返回。


class Solution:
    def rangeSum(self, nums: List[int], n: int, left: int, right: int) -> int:
        """二分查找+双指针求解前k小子数组和的和。"""
        ...

