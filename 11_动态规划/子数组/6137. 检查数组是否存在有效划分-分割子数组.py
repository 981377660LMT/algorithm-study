"""
刻子与顺子
如果获得的这些子数组中每个都能满足下述条件 之一 ，则可以称其为数组的一种 有效 划分：

子数组 恰 由 2 个相等元素组成，例如，子数组 [2,2] 。
子数组 恰 由 3 个相等元素组成，例如，子数组 [4,4,4] 。
子数组 恰 由 3 个连续递增元素组成，并且相邻元素之间的差值为 1 。例如，子数组 [3,4,5] ，但是子数组 [1,3,5] 不符合要求。

"""

from functools import lru_cache
from typing import List

# !没有意识到这个是dp hhh
# !用groupby做 发现 `groupby无法处理多种情况`
# !分割子数组：注意到答案只与当前位置(剩下的子数组)有关 所以是线性dp


class Solution:
    def validPartition(self, nums: List[int]) -> bool:
        """dp分割子数组"""
        n = len(nums)

        @lru_cache(None)
        def dfs(index: int) -> bool:
            if index >= n:
                return index == n

            res = False
            if index + 1 < n and nums[index] == nums[index + 1]:
                res |= dfs(index + 2)

            if index + 2 < n and (
                nums[index] == nums[index + 1] == nums[index + 2]
                or (nums[index] == (nums[index + 1] - 1) == (nums[index + 2] - 2))
            ):
                res |= dfs(index + 3)

            return res

        res = dfs(0)
        dfs.cache_clear()
        return res

    def validPartition2(self, nums: List[int]) -> bool:
        """dp分割子数组"""
        n = len(nums)
        dp = [False] * (n + 1)
        dp[0] = True
        for i in range(n):
            if i + 1 < n and nums[i] == nums[i + 1]:
                dp[i + 2] |= dp[i]
            if i + 2 < n and (
                nums[i] == nums[i + 1] == nums[i + 2]
                or (nums[i] == (nums[i + 1] - 1) == (nums[i + 2] - 2))
            ):
                dp[i + 3] |= dp[i]

        return dp[-1]


# print(Solution().validPartition(nums=[4, 4, 4, 5, 6]))
# print(Solution().validPartition(nums=[4, 4, 4]))
print(Solution().validPartition2(nums=[4, 4, 4, 5, 6, 7, 8, 9]))
# print(Solution().validPartition(nums=[1, 1, 1, 2]))
