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
        n = len(nums)
        if n == 1:
            return nums[0]

        dp = [[[0, 0] for _ in range(2)] for _ in range(n)]  # (index,pre,root) [不选 选]
        dp[0][1][1] = nums[0]
        for i in range(1, n):
            for pre in range(2):
                for cur in range(2):
                    if pre == cur == 1:
                        continue
                    for root in range(2):
                        dp[i][cur][root] = max(
                            dp[i][cur][root], dp[i - 1][pre][root] + (cur and nums[i])
                        )

        res = -INF
        for pre in range(2):
            for root in range(2):
                if pre == root == 1:
                    continue
                res = max(res, dp[n - 1][pre][root])
        return res
