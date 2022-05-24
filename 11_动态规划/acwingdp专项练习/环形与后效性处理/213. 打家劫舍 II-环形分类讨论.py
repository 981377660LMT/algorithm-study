# 偷最后一个/不偷最后一个


from functools import lru_cache
from typing import List


class Solution:
    def rob(self, nums: List[int]) -> int:
        @lru_cache(None)
        def dfs1(index: int, hasPre: int) -> int:
            """最后一个选了"""
            if index == n - 1:
                return nums[index] if not hasPre else -int(1e20)
            res = dfs1(index + 1, False)
            if not hasPre:
                res = max(res, dfs1(index + 1, True) + nums[index])
            return res

        @lru_cache(None)
        def dfs2(index: int, hasPre: int) -> int:
            """最后一个没选"""
            if index == n - 1:
                return 0
            res = dfs2(index + 1, False)
            if not hasPre:
                res = max(res, dfs2(index + 1, True) + nums[index])
            return res

        n = len(nums)
        if n == 1:  # 注意这里
            return nums[0]
        return max(dfs1(0, True), dfs2(0, False))

