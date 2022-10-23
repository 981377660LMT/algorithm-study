from bisect import bisect_right
from typing import List

# 11_动态规划/出租车问题/1235. 规划兼职工作.py


class Solution:
    def maxTaxiEarnings(self, n: int, rides: List[List[int]]) -> int:
        n = len(rides)
        rides.sort(key=lambda x: x[1])
        dp = [0] * (n + 1)

        for i in range(n):
            dp[i + 1] = dp[i]  # 不选
            start, end, tip = rides[i]  # 选
            score = end - start + tip
            prePos = bisect_right(rides, start, key=lambda x: x[1]) - 1
            dp[i + 1] = max(dp[i + 1], score + dp[prePos + 1])
        return dp[-1]


print(Solution().maxTaxiEarnings(n=5, rides=[[2, 5, 4], [1, 5, 1]]))
