# https://leetcode.cn/problems/minimum-number-of-coins-for-fruits-ii/description/
# 你在一个水果超市里，货架上摆满了玲琅满目的奇珍异果。
# 给你一个下标从 1 开始的数组 prices ，其中 prices[i] 表示你购买第 i 个水果需要花费的金币数目。
# 水果超市有如下促销活动：
# 如果你花费 price[i] 购买了水果 i ，那么接下来的 i 个水果你都可以免费获得。
# 注意 ，即使你 可以 免费获得水果 j ，你仍然可以花费 prices[j] 个金币去购买它以便能免费获得接下来的 j 个水果。
# 请你返回获得所有水果所需要的 最少 金币数。
#
# dp[i]:获得前i个水果`且最后一个水果i是花钱买的`，所需的最少金币数(0<=i<=n)
# dp[0]=0
# dp[i]=min(dp[j]+prices[i-1]) | i//2<=j<i)
# !维护窗口内的`dp[j]`的最小值即可
# !答案为min(dp[(n+1)//2],dp[n//2+1],...,dp[n])

from MonoQueue import MonoQueue

from typing import List, Tuple

INF = int(1e18)


class Solution:
    def minimumCoins(self, prices: List[int]) -> int:
        n = len(prices)
        queue = MonoQueue[Tuple[int, int]](lambda x, y: x[0] < y[0])  # (value, index)
        dp = [INF] * (n + 1)
        dp[0] = 0
        queue.append((dp[0], 0))
        for i in range(1, n + 1):
            while queue and queue.head()[1] < i // 2:
                queue.popleft()
            preMin = queue.head()[0] if queue else INF
            dp[i] = min(dp[i], preMin + prices[i - 1])
            queue.append((dp[i], i))
        return min(dp[(n + 1) // 2 : n + 1])


if __name__ == "__main__":
    # prices = [3,1,2]
    print(Solution().minimumCoins([3, 1, 2]))
    # prices = [1,10,1,1]
    print(Solution().minimumCoins([1, 10, 1, 1]))
    # [1,37,19,38,11,42,18,33,6,37,15,48,23,12,41,18,27,32]
    print(
        Solution().minimumCoins(
            [1, 37, 19, 38, 11, 42, 18, 33, 6, 37, 15, 48, 23, 12, 41, 18, 27, 32]
        )
    )
