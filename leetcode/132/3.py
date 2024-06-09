from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


def max2(a: int, b: int) -> int:
    return a if a > b else b


# 相邻不等元素的对数<=k，求最长子序列.
class Solution:
    def maximumLength(self, nums: List[int], k: int) -> int:
        @lru_cache(None)
        def dfs(index: int, pre: int, count: int) -> int:
            if count > k:
                return -INF
            if index == n:
                return 0

            res = 0
            cur = nums[index]
            bad = pre != -1 and cur != pre
            # 选
            res = max2(res, dfs(index + 1, cur, count + bad) + 1)
            # 不选
            res = max2(res, dfs(index + 1, pre, count))
            return res

        n = len(nums)
        res = dfs(0, -1, 0)
        dfs.cache_clear()
        return res


print(Solution().maximumLength([1, 2, 3, 4, 5, 1], 0))
