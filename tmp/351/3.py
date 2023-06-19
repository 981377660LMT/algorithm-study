from functools import lru_cache
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)
# 给你一个下标从 0 开始的整数数组 nums ，它包含 n 个 互不相同 的正整数。如果 nums 的一个排列满足以下条件，我们称它是一个特别的排列：

# 对于 0 <= i < n - 1 的下标 i ，要么 nums[i] % nums[i+1] == 0 ，要么 nums[i+1] % nums[i] == 0 。
# 请你返回特别排列的总数目，由于答案可能很大，请将它对 109 + 7 取余 后返回。


class Solution:
    def specialPerm(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs(index: int, pre: int, visited: int) -> int:
            if index == n:
                return 1
            res = 0
            for next in range(n):
                if visited & (1 << next) or (
                    pre != -1 and nums[next] % nums[pre] != 0 and nums[pre] % nums[next] != 0
                ):
                    continue
                res += dfs(index + 1, next, visited | 1 << next)
            return res % MOD

        n = len(nums)
        res = dfs(0, -1, 0)
        dfs.cache_clear()
        return res % MOD
