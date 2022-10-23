"""
983. 最低票价 
AutoX-3. 出行的最少购票费用
https://leetcode.cn/contest/autox2023/problems/BjAFy9/
"""
from bisect import bisect_left
from functools import lru_cache
from typing import List

# 火车票有k种不同的销售方式：
# tickets[i] = [duration_i, price_i]，表示第 i 种套票的有效天数和价格

# 返回你想要完成在给定的列表 days 中列出的每一天的旅行所需要的最低消费。
# !1 <= days.length <= 10^5
# !1 <= days[i] < days[i+1] <= 10^9
# !1 <= tickets.length <= 20

# !dp转移时需要二分查找


INF = int(1e18)


class Solution:
    def minCostToTravelOnDays(self, days: List[int], tickets: List[List[int]]) -> int:
        """https://leetcode.cn/contest/autox2023/problems/BjAFy9/"""
        n = len(days)
        dp = [INF] * (n + 1)  # !dp[i]表示时间到达前i天时的最少花费
        dp[0] = 0
        for i in range(n):
            cur = days[i]
            for retain, price in tickets:
                pre = cur - retain + 1
                pos = bisect_left(days, pre)
                dp[i + 1] = min(dp[i + 1], dp[pos] + price)
        return dp[-1]

    #########################################################################################################
    def mincostTickets(self, days: List[int], costs: List[int]) -> int:
        """https://leetcode.cn/problems/minimum-cost-for-tickets/"""
        n = days[-1]
        dSet = set(days)
        dp = [0] * (n + 1)

        for d in range(n + 1):
            if d not in dSet:  # 如果不在就继承前一个day的状态
                dp[d] = dp[d - 1]
            else:
                dp[d] = min(
                    dp[max(0, d - 1)] + costs[0],
                    dp[max(0, d - 7)] + costs[1],
                    dp[max(0, d - 30)] + costs[2],
                )

        return dp[-1]

    def mincostTickets2(self, days: List[int], costs: List[int]) -> int:
        """https://leetcode.cn/problems/minimum-cost-for-tickets/"""

        @lru_cache(None)
        def dfs(index: int) -> int:
            if index > n:
                return 0
            if index not in need:
                pos = bisect_left(days, index)
                next = days[pos] if pos < n else n + 1
                return dfs(next)
                return dfs(index + 1)  # !如果days很大 这里可以二分查找下一个位置

            return min(
                dfs(index + 1) + costs[0], dfs(index + 7) + costs[1], dfs(index + 30) + costs[2]
            )

        need = set(days)
        n = days[-1]
        res = dfs(days[0])
        dfs.cache_clear()
        return res


assert Solution().minCostToTravelOnDays([1, 2, 3, 4], [[1, 3], [2, 5], [3, 7]]) == 10
assert Solution().minCostToTravelOnDays([1, 4, 5], [[1, 4], [5, 6], [2, 5]]) == 6
#######################################################################
print(Solution().mincostTickets([1, 4, 6, 7, 8, 20], [2, 7, 15]))
print(Solution().mincostTickets2([1, 4, 6, 7, 8, 20], [2, 7, 15]))
# 输入：days = [1,4,6,7,8,20], costs = [2,7,15]
# 输出：11
# 解释：
# 例如，这里有一种购买通行证的方法，可以让你完成你的旅行计划：
# 在第 1 天，你花了 costs[0] = $2 买了一张为期 1 天的通行证，它将在第 1 天生效。
# 在第 3 天，你花了 costs[1] = $7 买了一张为期 7 天的通行证，它将在第 3, 4, ..., 9 天生效。
# 在第 20 天，你花了 costs[0] = $2 买了一张为期 1 天的通行证，它将在第 20 天生效。
# 你总共花了 $11，并完成了你计划的每一天旅行。
