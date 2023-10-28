from functools import lru_cache
from typing import List

INF = int(1e18)


class Solution:
    def rob(self, nums: List[int]) -> int:
        """dfs 考虑第一个选还是不选"""

        @lru_cache(None)
        def dfs(index: int, hasPre: bool, root: bool) -> int:
            """当前在index 前一个点是否选择 第一个点是否选择"""
            if index == n:
                return -INF if (hasPre and root) else 0  # !选了第一个，最后一个不能选
            res = dfs(index + 1, False, root)
            if not hasPre:
                res = max(res, dfs(index + 1, True) + nums[index], root)
            return res

        n = len(nums)
        if n == 1:  # !注意这里不成环
            return nums[0]
        return max(dfs(1, True, True) + nums[0], dfs(1, False, False))

    def rob2(self, nums: List[int]) -> int:
        """dp 考虑第一个选还是不选"""
        if not nums:
            return 0

        def cal0() -> int:  # 不偷第一个(i的范围是[1, n-1])
            dp0, dp1 = 0, 0
            for i in range(1, len(nums)):
                dp0, dp1 = max(dp0, dp1), max(dp0 + nums[i], dp1)
            return max(dp0, dp1)

        def cal1() -> int:  # 偷第一个(i的范围是[2, n-2])
            dp0, dp1 = 0, 0
            for i in range(2, len(nums) - 1):
                dp0, dp1 = max(dp0, dp1), max(dp0 + nums[i], dp1)
            return max(dp0, dp1)

        return max(cal0(), cal1() + nums[0])
