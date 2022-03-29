from functools import lru_cache
from typing import List


class Solution:
    def maxSumTwoNoOverlap(self, nums: List[int], k: int):
        """
        k个非重叠子数组的最大和
        子数组：取或全不取
        """

        @lru_cache(None)
        def dfs(index: int, remain: int, isPreSelected: bool) -> int:
            if remain < 0:
                return -int(1e20)
            if index == n:
                return 0 if remain == 0 else -int(1e20)

            res = -int(1e20)
            skip = dfs(index + 1, remain, False)
            choose1 = nums[index] + dfs(index + 1, remain - 1, True)
            if isPreSelected:
                choose2 = nums[index] + dfs(index + 1, remain, True)
                res = max(res, choose2)

            res = max(res, choose1, skip)
            return res

        n = len(nums)
        return dfs(0, k, False)


print(Solution().maxSumTwoNoOverlap(nums=[10, -2, 1, 0, 5, -25, 10, -10, 5], k=3))
