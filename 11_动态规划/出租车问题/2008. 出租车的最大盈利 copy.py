from bisect import bisect_right
from typing import List


class Solution:
    def maxTaxiEarnings(self, n: int, rides: List[List[int]]) -> int:
        rides.sort(key=lambda x: x[1])
        dp = [e - s + t for s, e, t in rides]
        ends = [e for _, e, _ in rides]

        for i in range(1, len(rides)):
            pre = bisect_right(ends, rides[i][0]) - 1
            if pre >= 0:
                dp[i] = max(dp[i - 1], dp[pre] + rides[i][1] - rides[i][0] + rides[i][2])
            else:
                dp[i] = max(dp[i - 1], rides[i][1] - rides[i][0] + rides[i][2])
        return dp[-1]


print(Solution().maxTaxiEarnings(n=5, rides=[[2, 5, 4], [1, 5, 1]]))
