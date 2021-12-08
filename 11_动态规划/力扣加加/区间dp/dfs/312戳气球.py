from typing import List
from functools import lru_cache

# 1 <= n <= 500
# 估计O(n^3)
class Solution:
    def maxCoins(self, nums: List[int]) -> int:
        nums = [1] + nums + [1]

        @lru_cache(None)
        def dfs(left: int, right: int) -> int:
            if left >= right - 1:
                return 0
            res = 0
            for mid in range(left + 1, right):
                res = max(
                    res, nums[left] * nums[mid] * nums[right] + dfs(left, mid) + dfs(mid, right)
                )

            return res

        return dfs(0, len(nums) - 1)


print(Solution().maxCoins(nums=[3, 1, 5, 8]))
# 输出：167
# 解释：
# nums = [3,1,5,8] --> [3,5,8] --> [3,8] --> [8] --> []
# coins =  3*1*5    +   3*5*8   +  1*3*8  + 1*8*1 = 167
# 戳破第 i 个气球，你可以获得 nums[i - 1] * nums[i] * nums[i + 1] 枚硬币。
