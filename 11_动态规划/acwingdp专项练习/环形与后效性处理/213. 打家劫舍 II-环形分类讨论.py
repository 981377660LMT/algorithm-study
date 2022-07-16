from functools import lru_cache
from typing import List


class Solution:
    def rob(self, nums: List[int]) -> int:
        """dfs 考虑最后一个选还是不选"""

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

    def rob2(self, nums: List[int]) -> int:
        """dp 考虑第一个选还是不选"""

        def cal1(nums: List[int]) -> int:
            """第一个不选"""
            n = len(nums)
            dp = [0] * n
            for i in range(1, n):
                dp[i] = max(dp[i - 1], nums[i] + dp[i - 2])
            return dp[-1]

        def cal2(nums: List[int]) -> int:
            """第一个选"""
            n = len(nums)
            dp = [0] * n
            dp[0] = nums[0]
            for i in range(1, n):
                # 那么最后一个就不能选了
                dp[i] = max(dp[i - 1], (nums[i] if i != n - 1 else 0) + dp[i - 2])
            return dp[-1]

        n = len(nums)
        if n == 1:
            return nums[0]
        if n == 2:
            return max(nums[0], nums[1])
        return max(cal1(nums), cal2(nums))
