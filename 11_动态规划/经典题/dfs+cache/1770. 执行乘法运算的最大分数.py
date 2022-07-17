from typing import List
from functools import lru_cache

# 给你两个长度分别 n 和 m 的整数数组 nums 和 multipliers ，其中 n >= m ，数组下标 从 1 开始 计数。
# 在第 i 步操作（从 1 开始 计数）中，需要：

# 选择数组 nums 开头处或者末尾处 的整数 x 。
# 你获得 multipliers[i] * x 分，并累加到你的分数中。
# 将 x 从数组 nums 中移除。
# 在执行 m 步操作后，返回 最大 分数。

# -1000 <= nums[i], multipliers[i] <= 1000
# !m<=1e3,n<=1e5

# !1. 手写 max 快了 1000ms
# !2. 不用cache_clear() 超时


class Solution:
    def maximumScore(self, nums: List[int], multipliers: List[int]) -> int:
        @lru_cache(None)
        def dfs(left: int, right: int, index: int) -> int:
            """注意index由left,right唯一决定 因此复杂度仍然是O(m^2)的"""
            if index == m:
                return 0

            cand1 = nums[left] * multipliers[index] + dfs(left + 1, right, index + 1)
            cand2 = nums[right] * multipliers[index] + dfs(left, right - 1, index + 1)
            if cand1 > cand2:
                return cand1
            return cand2

        n, m = len(nums), len(multipliers)
        res = dfs(0, n - 1, 0)
        dfs.cache_clear()
        return res


print(Solution().maximumScore(nums=[1, 2, 3], multipliers=[3, 2, 1]))
print(Solution().maximumScore(nums=[2, 5, 4, 3, 1], multipliers=[3, 5, 1, 2, 4]))
# 输出：14
# 解释：一种最优解决方案如下：
# - 选择末尾处的整数 3 ，[1,2,3] ，得 3 * 3 = 9 分，累加到分数中。
# - 选择末尾处的整数 2 ，[1,2] ，得 2 * 2 = 4 分，累加到分数中。
# - 选择末尾处的整数 1 ，[1] ，得 1 * 1 = 1 分，累加到分数中。
# 总分数为 9 + 4 + 1 = 14 。
