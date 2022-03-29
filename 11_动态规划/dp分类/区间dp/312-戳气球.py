# 戳破第 i 个气球，你可以获得 nums[i - 1] * nums[i] * nums[i + 1] 枚硬币。
# You can take any candy at index i and you will get candies[i - 1] * candies[i] * candies[i + 1] candies
from functools import lru_cache

# n<=200
class Solution:
    def solve(self, candies):
        @lru_cache(None)
        def dfs(i, j):
            res = 0
            for k in range(i + 1, j):
                res = max(res, dfs(i, k) + dfs(k, j) + candies[i] * candies[k] * candies[j])
            return res

        candies = [1] + candies + [1]
        return dfs(0, len(candies) - 1)


print(Solution().solve([3, 1, 5, 8]))

