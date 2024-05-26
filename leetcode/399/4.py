from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums 和一个二维数组 queries，其中 queries[i] = [posi, xi]。

# 对于每个查询 i，首先将 nums[posi] 设置为 xi，然后计算查询 i 的答案，该答案为 nums 中 不包含相邻元素 的子序列的 最大 和。

# 返回所有查询的答案之和。

# 由于最终答案可能非常大，返回其对 109 + 7 取余 的结果。


# 子序列 是指从另一个数组中删除一些或不删除元素而不改变剩余元素顺序得到的数组。
class Solution:
    def maximumSumSubsequence(self, nums: List[int], queries: List[List[int]]) -> int:
        ...
