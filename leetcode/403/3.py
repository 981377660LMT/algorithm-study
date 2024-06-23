from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个长度为 n 的整数数组 nums。

# 子数组 nums[l..r]（其中 0 <= l <= r < n）的 成本 定义为：

# cost(l, r) = nums[l] - nums[l + 1] + ... + nums[r] * (−1)r − l

# 你的任务是将 nums 分割成若干子数组，使得所有子数组的成本之和 最大化，并确保每个元素 正好 属于一个子数组。

# 具体来说，如果 nums 被分割成 k 个子数组，且分割点为索引 i1, i2, ..., ik − 1（其中 0 <= i1 < i2 < ... < ik - 1 < n - 1），则总成本为：

# cost(0, i1) + cost(i1 + 1, i2) + ... + cost(ik − 1 + 1, n − 1)

# 返回在最优分割方式下的子数组成本之和的最大值。


# 注意：如果 nums 没有被分割，即 k = 1，则总成本即为 cost(0, n - 1)。


# 这个数负时，前一个数必须为正


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maximumTotalCost(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, prePlus: bool) -> int:
            if index == n:
                return 0
            cur = nums[index]
            if prePlus:
                return max2(dfs(index + 1, True) + cur, dfs(index + 1, False) - cur)
            else:
                return dfs(index + 1, True) + cur

        n = len(nums)
        res = dfs(0, False)
        dfs.cache_clear()
        return res
