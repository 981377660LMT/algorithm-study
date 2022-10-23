# 出行的最少购票费用
# 火车票有k种不同的销售方式：
# tickets[i] = [duration_i, price_i]，表示第 i 种套票的`有效天数和价格`
# !例如：乘客购买了有效天数为 n 的套票，则该套票在第 date ~ date+n-1 天期间都可以使用。
# 返回你想要完成在给定的列表 days 中列出的每一天的旅行所需要的最低消费。
# !1 <= days.length <= 10^5
# !1 <= days[i] < days[i+1] <= 10^9
# !1 <= tickets.length <= 20

# !dp转移时需要二分查找

from bisect import bisect_left
from typing import List

INF = int(1e18)


class Solution:
    def minCostToTravelOnDays(self, days: List[int], tickets: List[List[int]]) -> int:
        n = len(days)
        dp = [INF] * (n + 1)  # !dp[i]表示days[i]作为结尾时的最少花费
        dp[0] = 0
        for i in range(n):
            curDay = days[i]
            for retain, price in tickets:
                preDay = curDay - retain + 1  # !在之前这天买票
                pos = bisect_left(days, preDay)
                dp[i + 1] = min(dp[i + 1], dp[pos] + price)
        return dp[-1]


assert Solution().minCostToTravelOnDays(days=[1, 2, 3, 4], tickets=[[1, 3], [2, 5], [3, 7]]) == 10
assert Solution().minCostToTravelOnDays(days=[1, 4, 5], tickets=[[1, 4], [5, 6], [2, 5]]) == 6
