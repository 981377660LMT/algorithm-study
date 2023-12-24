from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始的 正 整数数组 nums 。

# 如果 nums 的一个子数组满足：移除这个子数组后剩余元素 严格递增 ，那么我们称这个子数组为 移除递增 子数组。比方说，[5, 3, 4, 6, 7] 中的 [3, 4] 是一个移除递增子数组，因为移除该子数组后，[5, 3, 4, 6, 7] 变为 [5, 6, 7] ，是严格递增的。

# 请你返回 nums 中 移除递增 子数组的总数目。

# 注意 ，剩余元素为空的数组也视为是递增的。


# 子数组 指的是一个数组中一段连续的元素序列。
class Solution:
    def incremovableSubarrayCount(self, nums: List[int]) -> int:
        #  枚举子数组
        res = 0
        n = len(nums)
        for i in range(n):
            for j in range(i, n):
                cur = nums[:i] + nums[j + 1 :]
                if all(cur[i] < cur[i + 1] for i in range(len(cur) - 1)):
                    res += 1
        return res
