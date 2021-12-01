from typing import List
from math import floor, ceil

# topk
class Solution:
    def minimizeError(self, prices: List[str], target: int) -> str:
        prices_ = [float(i) for i in prices]

        floor_prices = [floor(i) for i in prices_]
        ceil_prices = [ceil(i) for i in prices_]
        Min, Max = sum(floor_prices), sum(ceil_prices)
        print(Min, Max, sum(prices_))
        if target < Min or target > Max:
            return '-1'

        k = target - Min  # 需要选k个数取ceil => 选取代价最小的那k个
        diff_sum = sum(prices_) - Min
        diff = [
            abs(floor_prices[i] - prices_[i]) - abs(ceil_prices[i] - prices_[i])
            for i in range(len(prices_))
        ]

        diff = sorted(diff, reverse=True)
        print(diff)

        return '%.3f' % (diff_sum - sum(diff[:k]))


print(Solution().minimizeError(prices=["0.700", "2.800", "4.900"], target=8))
# 输出："1.000"
# 解释：
# 使用 Floor，Ceil 和 Ceil 操作得到 (0.7 - 0) + (3 - 2.8) + (5 - 4.9) = 0.7 + 0.2 + 0.1 = 1.0 。

