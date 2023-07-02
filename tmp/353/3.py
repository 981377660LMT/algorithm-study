from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始的整数数组 nums 。nums 的一个子数组如果满足以下条件，那么它是 不间断 的：

# i，i + 1 ，...，j  表示子数组中的下标。对于所有满足 i <= i1, i2 <= j 的下标对，都有 0 <= |nums[i1] - nums[i2]| <= 2 。
# 请你返回 不间断 子数组的总数目。

# 子数组是一个数组中一段连续 非空 的元素序列。


class Solution:
    def continuousSubarrays(self, nums: List[int]) -> int:
        res, left, n = 0, 0, len(nums)
        sl = SortedList()
        for right in range(n):
            sl.add(nums[right])
            while left <= right and sl[-1] - sl[0] > 2:
                sl.remove(nums[left])
                left += 1
            res += right - left + 1
        return res
