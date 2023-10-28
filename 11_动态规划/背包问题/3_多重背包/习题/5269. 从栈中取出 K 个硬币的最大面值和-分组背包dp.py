from itertools import accumulate
from typing import List

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= k <= sum(piles[i].length) <= 2000
# 1 <= n <= 1000

# 这道题时间复杂度为O(k*sum(piles[i].length))


class Solution:
    def maxValueOfCoins(self, piles: List[List[int]], k: int) -> int:
        """时间复杂度O(背包容量*物品个数)"""
        n = len(piles)
        preSums = []
        for i in range(n):
            preSums.append([0] + list(accumulate(piles[i])))

        dp = [0] * (k + 1)
        for i in range(n):
            for j in range(k, -1, -1):
                for select in range(len(piles[i]) + 1):
                    sum_ = preSums[i][select]
                    if j - select >= 0:
                        dp[j] = max(dp[j], dp[j - select] + sum_)

        return dp[-1]


print(Solution().maxValueOfCoins(piles=[[1, 100, 3], [7, 8, 9]], k=2))
print(
    Solution().maxValueOfCoins(
        piles=[[100], [100], [100], [100], [100], [100], [1, 1, 1, 1, 1, 1, 700]], k=7
    )
)
