from typing import List
from functools import lru_cache

# 在第 i 步操作（从 1 开始 计数）中，需要：

# 选择数组 nums 开头处或者末尾处 的整数 x 。
# 你获得 multipliers[i] * x 分，并累加到你的分数中。
# 将 x 从数组 nums 中移除。
# 在执行 m 步操作后，返回 最大 分数。

# -1000 <= nums[i], multipliers[i] <= 1000
# m <= n <= 105


class Solution:
    def maximumScore(self, nums: List[int], multipliers: List[int]) -> int:
        n, m = len(nums), len(multipliers)

        @lru_cache(None)
        def dfs(l: int, r: int, index: int) -> int:
            if index == m:
                return 0

            return max(
                nums[l] * multipliers[index] + dfs(l + 1, r, index + 1),
                nums[r] * multipliers[index] + dfs(l, r - 1, index + 1),
            )

        res = dfs(0, n - 1, 0)
        # print(dfs.cache_info())
        dfs.cache_clear()

        return res


print(Solution().maximumScore(nums=[1, 2, 3], multipliers=[3, 2, 1]))
# 输出：14
# 解释：一种最优解决方案如下：
# - 选择末尾处的整数 3 ，[1,2,3] ，得 3 * 3 = 9 分，累加到分数中。
# - 选择末尾处的整数 2 ，[1,2] ，得 2 * 2 = 4 分，累加到分数中。
# - 选择末尾处的整数 1 ，[1] ，得 1 * 1 = 1 分，累加到分数中。
# 总分数为 9 + 4 + 1 = 14 。

