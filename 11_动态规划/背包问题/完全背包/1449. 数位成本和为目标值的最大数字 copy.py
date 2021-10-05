from typing import List


class Solution:
    def largestNumber(self, cost: List[int], target: int) -> str:
        dp = [''] * (target + 1)
        for i in range(9):
            for j in range(cost[i], target + 1):
                pre = dp[j - cost[i]]
                if j!=cost[i] and len(pre)
        return dp[target] if dp[target] else '0'


print(Solution().largestNumber([4, 3, 2, 5, 6, 7, 2, 5, 5], 9))
