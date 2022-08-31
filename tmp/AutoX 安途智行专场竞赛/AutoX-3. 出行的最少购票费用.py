from bisect import bisect_left
from typing import List

INF = int(1e18)


class Solution:
    def minCostToTravelOnDays(self, days: List[int], tickets: List[List[int]]) -> int:
        n = len(days)
        dp = [INF] * (n + 1)
        dp[0] = 0
        for i in range(1, n + 1):
            cur = days[i - 1]
            for dur, cost in tickets:
                pre = bisect_left(days, cur - dur + 1)
                dp[i] = min(dp[i], dp[pre] + cost)
        return dp[n]


print(Solution().minCostToTravelOnDays(days=[1, 2, 3, 4], tickets=[[1, 3], [2, 5], [3, 7]]))
