from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个下标从 0 开始、长度为 n 的整数数组 nums ，和一个整数 k 。

# 你可以执行下述 递增 运算 任意 次（可以是 0 次）：

# 从范围 [0, n - 1] 中选则一个下标 i ，并将 nums[i] 的值加 1 。
# 如果数组中任何长度 大于或等于 3 的子数组，其 最大 元素都大于或等于 k ，则认为数组是一个 美丽数组 。

# 以整数形式返回使数组变为 美丽数组 需要执行的 最小 递增运算数。


# 子数组是数组中的一个连续 非空 元素序列。


def min(a, b):
    return a if a < b else b


def max(a, b):
    return a if a > b else b


class Solution:
    def minIncrementOperations(self, nums: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, ok1: bool, ok2: bool) -> int:
            if index == n:
                return 0

            cur = nums[index]
            # 当前加
            res = dfs(index + 1, ok2, True) + max(0, k - cur)

            # 当前不加
            curOk = cur >= k
            if ok1 or ok2 or curOk:
                res = min(res, dfs(index + 1, ok2, curOk))

            return res

        n = len(nums)
        res = dfs(0, True, True)
        dfs.cache_clear()
        return res
