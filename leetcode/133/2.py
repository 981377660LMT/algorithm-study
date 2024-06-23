from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二进制数组 nums 。

# 你可以对数组执行以下操作 任意 次（也可以 0 次）：

# 选择数组中 任意连续 3 个元素，并将它们 全部反转 。
# 反转 一个元素指的是将它的值从 0 变 1 ，或者从 1 变 0 。


# 请你返回将 nums 中所有元素变为 1 的 最少 操作次数。如果无法全部变成 1 ，返回 -1
class Solution:
    def minOperations(self, nums: List[int]) -> int:
        res = 0
        for i in range(len(nums)):
            if nums[i] == 0:
                if i + 2 >= len(nums):
                    return -1
                nums[i] ^= 1
                nums[i + 1] ^= 1
                nums[i + 2] ^= 1
                res += 1
        return res
