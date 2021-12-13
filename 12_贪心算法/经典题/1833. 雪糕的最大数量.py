from typing import List

# 1 <= n <= 105
# 1 <= costs[i] <= 105
class Solution:
    def maxIceCream(self, costs: List[int], coins: int) -> int:
        costs = sorted(costs)
        for i, v in enumerate(costs):
            coins -= v
            if coins < 0:
                return i
        return len(costs)


print(Solution().maxIceCream(costs=[1, 3, 2, 4, 1], coins=7))
# 输出：4
# 解释：Tony 可以买下标为 0、1、2、4 的雪糕，总价为 1 + 3 + 2 + 1 = 7
