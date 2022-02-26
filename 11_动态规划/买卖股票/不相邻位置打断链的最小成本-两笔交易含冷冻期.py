# 给出了一个由N个整数组成的数组A。数组A中的元素合起来表示一个链，每个元素表示链中每个环节的强度。
# 我们要把这条链分成三个小点的链。
# 我们所能做的就是在两个不相邻的位置打断这个链。

# 买卖股票的最佳时机含冷冻期
from typing import List

INF = int(1e20)


class Solution:
    def minCost(self, prices: List[int]) -> int:
        # dp[i][0] 分隔了零次
        # dp[i][1] 分隔了一次,当天能选
        # dp[i][2] 分隔了一次,当天不能选
        # dp[i][3] 分隔了两次

        prices = prices[1:-1]
        dp = [[INF] * 4 for _ in range(len(prices))]
        dp[0][0] = 0
        dp[0][1] = prices[0]
        dp[0][2] = INF
        dp[0][3] = INF

        for i in range(1, len(prices)):
            dp[i][0] = 0
            dp[i][1] = min(dp[i - 1][1], prices[i])
            dp[i][2] = dp[i - 1][1]
            dp[i][3] = min(dp[i - 1][1] + prices[i], dp[i - 1][2] + prices[i], dp[i - 1][3])

        return dp[-1][3]


print(Solution().minCost([5, 2, 4, 6, 3, 7]))  # 5

