from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的 非递减 整数数组 nums 。

# 你可以执行以下操作任意次：


# 选择 两个 下标 i 和 j ，满足 i < j 且 nums[i] < nums[j] 。
# 将 nums 中下标在 i 和 j 处的元素删除。剩余元素按照原来的顺序组成新的数组，下标也重新从 0 开始编号。
# 请你返回一个整数，表示执行以上操作任意次后（可以执行 0 次），nums 数组的 最小 数组长度。


class Solution:
    def minLengthAfterRemovals(self, nums: List[int]) -> int:
        if len(nums) == 1:
            return 1
        counter = Counter(nums)
        maxFreq = max(counter.values())
        otherFreq = len(nums) - maxFreq
        rest = maxFreq - otherFreq
        if rest <= 0:
            return 1 if len(nums) & 1 else 0
        return rest
