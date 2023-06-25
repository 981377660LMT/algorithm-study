from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个二元数组 nums 。

# 如果数组中的某个子数组 恰好 只存在 一 个值为 1 的元素，则认为该子数组是一个 好子数组 。

# 请你统计将数组 nums 划分成若干 好子数组 的方法数，并以整数形式返回。由于数字可能很大，返回其对 109 + 7 取余 之后的结果。


# 子数组是数组中的一个连续 非空 元素序列。
class Solution:
    def numberOfGoodSubarraySplits(self, nums: List[int]) -> int:
        if 1 not in nums:
            return 0

        @lru_cache(None)
        def dfs(index: int, curOne: int) -> int:
            if index == n:
                return int(curOne == 1)
            res = 0
            cur = nums[index]
            if curOne == 0:
                if cur == 1:
                    res += dfs(index + 1, 1)
                else:
                    res += dfs(index + 1, 0)
            else:
                if cur == 1:
                    res += dfs(index + 1, 1)
                else:
                    res += dfs(index + 1, 1)
                    res += dfs(index + 1, 0)

            return res % MOD

        n = len(nums)
        res = dfs(0, 0)
        dfs.cache_clear()
        return res % MOD
