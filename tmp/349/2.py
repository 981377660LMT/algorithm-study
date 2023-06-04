from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始、长度为 n 的整数排列 nums 。

# 如果排列的第一个数字等于 1 且最后一个数字等于 n ，则称其为 半有序排列 。你可以执行多次下述操作，直到将 nums 变成一个 半有序排列 ：

# 选择 nums 中相邻的两个元素，然后交换它们。
# 返回使 nums 变成 半有序排列 所需的最小操作次数。

# 排列 是一个长度为 n 的整数序列，其中包含从 1 到 n 的每个数字恰好一次。


class Solution:
    def semiOrderedPermutation(self, nums: List[int]) -> int:
        n = len(nums)
        index1 = nums.index(1)
        indexn = nums.index(n)
        minus = index1 > indexn
        return (index1 - 0) + (n - 1 - indexn) - minus
