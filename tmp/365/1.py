from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的整数数组 nums 。

# 请你从所有满足 i < j < k 的下标三元组 (i, j, k) 中，找出并返回下标三元组的最大值。如果所有满足条件的三元组的值都是负数，则返回 0 。


# 下标三元组 (i, j, k) 的值等于 (nums[i] - nums[j]) * nums[k] 。
class Solution:
    def maximumTripletValue(self, nums: List[int]) -> int:
        # 枚举j
        left = SortedList()
        right = SortedList(nums[1:])
        res = 0
        for j in range(1, len(nums) - 1):
            left.add(nums[j - 1])
            right.remove(nums[j])
            res = max(res, (left[-1] - nums[j]) * right[-1])

        return res


# nums = [12,6,1,2,7]

print(Solution().maximumTripletValue(nums=[12, 6, 1, 2, 7]))
