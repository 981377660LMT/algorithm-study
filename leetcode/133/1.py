from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 。一次操作中，你可以将 nums 中的 任意 一个元素增加或者减少 1 。


# 请你返回将 nums 中所有元素都可以被 3 整除的 最少 操作次数。
class Solution:
    def minimumOperations(self, nums: List[int]) -> int:
        res = 0
        for i in nums:
            mod_ = i % 3
            if mod_ == 0:
                continue
            if mod_ == 1:
                res += 1
            if mod_ == 2:
                res += 1
        return res
