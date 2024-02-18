from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个整数数组 nums ，如果 nums 至少 包含 2 个元素，你可以执行以下操作中的 任意 一个：

# 选择 nums 中最前面两个元素并且删除它们。
# 选择 nums 中最后两个元素并且删除它们。
# 选择 nums 中第一个和最后一个元素并且删除它们。
# 一次操作的 分数 是被删除元素的和。

# 在确保 所有操作分数相同 的前提下，请你求出 最多 能进行多少次操作。


# 请你返回按照上述要求 最多 可以进行的操作次数。


def max2(a: int, b: int) -> int:
    return a if a > b else b


class Solution:
    def maxOperations(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(start: int, end: int, target: int) -> int:
            if end - start <= 1:
                return 0
            res = 0
            if nums[start] + nums[start + 1] == target:
                res = max2(res, dfs(start + 2, end, target) + 1)
            if nums[start] + nums[end - 1] == target:
                res = max2(res, dfs(start + 1, end - 1, target) + 1)
            if nums[end - 1] + nums[end - 2] == target:
                res = max2(res, dfs(start, end - 2, target) + 1)
            return res

        n = len(nums)
        res1 = dfs(2, n, nums[0] + nums[1]) + 1
        res2 = dfs(1, n - 1, nums[0] + nums[n - 1]) + 1
        res3 = dfs(0, n - 2, nums[n - 1] + nums[n - 2]) + 1
        dfs.cache_clear()
        return max2(res1, max2(res2, res3))
