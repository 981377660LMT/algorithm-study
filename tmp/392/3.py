from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 和一个 非负 整数 k 。

# 一次操作中，你可以选择任一下标 i ，然后将 nums[i] 加 1 或者减 1 。

# 请你返回将 nums 中位数 变为 k 所需要的 最少 操作次数。


# 一个数组的 中位数 指的是数组按 非递减 顺序排序后最中间的元素。如果数组长度为偶数，我们选择中间两个数的较大值为中位数。


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def minOperationsToMakeMedianK(self, nums: List[int], k: int) -> int:
        nums.sort()
        n = len(nums)
        mid = n // 2
        curSum = 0
        for i, v in enumerate(nums):
            if i < mid:
                curSum += max2(0, v - k)
            elif i == mid:
                curSum += abs(k - v)
            else:
                curSum += max2(0, k - v)
        return curSum
