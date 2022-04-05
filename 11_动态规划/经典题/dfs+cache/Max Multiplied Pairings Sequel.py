# 首尾选择数字相乘 求最大值
from functools import lru_cache


class Solution:
    def solve(self, nums, multipliers):
        @lru_cache(None)
        def dfs(left, right, index):
            # There are n choices of left and n choices of right corresponding to each left. Thus n^2
            if index == len(multipliers):
                return 0

            return max(
                dfs(left + 1, right, index + 1) + multipliers[index] * nums[left],
                dfs(left, right - 1, index + 1) + multipliers[index] * nums[right],
            )

        return dfs(0, len(nums) - 1, 0)


print(Solution().solve(nums=[5, 2, -7], multipliers=[2, 4, -1]))
# We can get 5 * 2 + 2 * 4 + -7 * -1 = 25
