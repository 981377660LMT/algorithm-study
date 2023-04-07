from typing import List


class Solution:
    def maxProfit(self, k: int, prices: List[int]) -> int:
        if not prices:
            return 0

        n = len(prices)
        # 二分查找的上下界
        left, right = 1, max(prices)
        # 存储答案，如果值为 -1 表示二分查找失败
        ans = -1

        while left <= right:
            # 二分得到当前的斜率（手续费）
            c = (left + right) // 2

            # 使用与 714 题相同的动态规划方法求解出最大收益以及对应的交易次数
            buyCount = sellCount = 0
            buy, sell = -prices[0], 0

            for i in range(1, n):
                if sell - prices[i] >= buy:
                    buy = sell - prices[i]
                    buyCount = sellCount
                if buy + prices[i] - c >= sell:
                    sell = buy + prices[i] - c
                    sellCount = buyCount + 1

            # 如果交易次数大于等于 k，那么可以更新答案
            # 这里即使交易次数严格大于 k，更新答案也没有关系，因为总能二分到等于 k 的
            if sellCount >= k:
                # 别忘了加上 kc
                ans = sell + k * c
                left = c + 1
                print(ans, sellCount, sell, c)
            else:
                right = c - 1
        # 如果二分查找失败，说明交易次数的限制不是瓶颈
        # 可以看作交易次数无限，直接使用贪心方法得到答案
        if ans == -1:
            ans = sum(max(prices[i] - prices[i - 1], 0) for i in range(1, n))

        return ans


# 2
# [3,3,5,0,0,3,1,4]

print(Solution().maxProfit(2, [3, 3, 5, 0, 0, 3, 1, 4]))
