from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 和一个整数 k 。

# 请你用整数形式返回 nums 中的特定元素之 和 ，这些特定元素满足：其对应下标的二进制表示中恰存在 k 个置位。

# 整数的二进制表示中的 1 就是这个整数的 置位 。


# 例如，21 的二进制表示为 10101 ，其中有 3 个置位。
class Solution:
    def sumIndicesWithKSetBits(self, nums: List[int], k: int) -> int:
        res = 0
        for i, num in enumerate(nums):
            if bin(i).count("1") == k:
                res += num
        return res
