# 给你一个长度为 n 的数组 nums 和一个整数 m 。
# 请你判断能否执行一系列操作，将数组拆分成 n 个 非空 数组。

# 在每一步操作中，你可以选择一个 长度至少为 2 的现有数组（之前步骤的结果） 并将其拆分成 2 个子数组，而得到的 每个 子数组，至少 需要满足以下条件之一：

# 子数组的长度为 1 ，或者
# 子数组元素之和 大于或等于  m 。
# 如果你可以将给定数组拆分成 n 个满足要求的数组，返回 true ；否则，返回 false 。
# 注意：子数组是数组中的一个连续非空元素序列。


# n<=100

from functools import lru_cache
from itertools import accumulate
from typing import List


class Solution:
    def canSplitArray(self, nums: List[int], m: int) -> bool:
        if len(nums) <= 2:
            return True

        # 区间dp
        @lru_cache(None)
        def dfs(left: int, right: int) -> bool:
            if left >= right:
                return False
            if right - left == 1:
                return True
            if preSum[right] - preSum[left] < m:
                return False
            for k in range(left + 1, right):
                if dfs(left, k) and dfs(k, right):
                    return True
            return False

        preSum = [0] + list(accumulate(nums))
        res = dfs(0, len(nums))
        dfs.cache_clear()
        return res


# [2, 1, 3]
print(Solution().canSplitArray(nums=[2, 1, 3], m=5))
