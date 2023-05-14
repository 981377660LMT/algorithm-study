from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始长度为 n 的整数数组 nums 和一个整数 k 。每一次操作中，你可以选择一个数并将它乘 2 。

# 你最多可以进行 k 次操作，请你返回 nums[0] | nums[1] | ... | nums[n - 1] 的最大值。


# a | b 表示两个整数 a 和 b 的 按位或 运算。


class Solution:
    def matrixSum(self, nums: List[List[int]]) -> int:
        for v in nums:
            v.sort()
        mat = [list(v) for v in zip(*nums)]
        return sum(max(v) for v in mat)
