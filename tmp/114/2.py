from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的正整数数组 nums 。

# 你可以对数组执行以下两种操作 任意次 ：


# 从数组中选择 两个 值 相等 的元素，并将它们从数组中 删除 。
# 从数组中选择 三个 值 相等 的元素，并将它们从数组中 删除 。
# 请你返回使数组为空的 最少 操作次数，如果无法达成，请返回 -1 。


# 2x+3y
class Solution:
    def minOperations(self, nums: List[int]) -> int:
        counter = Counter(nums)
        res = 0
        for v in counter.values():
            if v == 1:
                return -1
            if v % 3 == 0:
                res += v // 3
            elif v % 3 == 1:  # 4, 7, 10
                res += (v - 4) // 3 + 2
            else:
                res += v // 3 + 1
        return res
