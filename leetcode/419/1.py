from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个由 n 个整数组成的数组 nums，以及两个整数 k 和 x。

# 数组的 x-sum 计算按照以下步骤进行：

# 统计数组中所有元素的出现次数。
# 仅保留出现次数最多的前 x 个元素的每次出现。如果两个元素的出现次数相同，则数值 较大 的元素被认为出现次数更多。
# 计算结果数组的和。
# 注意，如果数组中的不同元素少于 x 个，则其 x-sum 是数组的元素总和。

# 返回一个长度为 n - k + 1 的整数数组 answer，其中 answer[i] 是 子数组 nums[i..i + k - 1] 的 x-sum。


# 子数组 是数组内的一个连续 非空 的元素序列。
class Solution:
    def findXSum(self, nums: List[int], k: int, x: int) -> List[int]:
        # def f(nums: List[int], x: int) -> int:
        #     counter = Counter(nums)
        #     sl = SortedList()
        #     for key, value in counter.items():
        #         sl.add((value, key))
        #     sl = sl[-x:]
        #     return sum([value * key for key, value in sl])

        # return [f(nums[i : i + k], x) for i in range(len(nums) - k + 1)]

        ...
