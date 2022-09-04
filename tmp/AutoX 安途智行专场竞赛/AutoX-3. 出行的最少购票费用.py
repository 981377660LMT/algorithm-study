# 出行的最少购票费用
# 航空公司向经常乘坐飞机的乘客们提供了一些商务套票，
# tickets[i] = [duration_i, price_i]，表示第 i 种套票的有效天数和价格。
# 乘客购买了有效天数为 n 的套票，则该套票在第 date ~ date+n-1 天期间都可以使用。
# 现有一名乘客将在未来的几天中出行，
# days[i] 表示他第 i 次出行的时间，如果他选择购买商务套票，请返回他将花费的最少金额。

# n<=1e5
# 1 <= tickets.length <= 20
# 1 <= days[i] < days[i+1] <= 10^9

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
