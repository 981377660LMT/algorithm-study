from typing import List

# 多少个子数组 单调递减相邻差为1
class Solution:
    def getDescentPeriods(self, prices: List[int]) -> int:
        n = len(prices)
        res = 1
        dp = 1
        for i in range(1, n):
            if prices[i] == prices[i - 1] - 1:
                dp += 1
            else:
                dp = 1
            res += dp
        return res


print(Solution().getDescentPeriods(prices=[3, 2, 1, 4]))
print(Solution().getDescentPeriods(prices=[8, 6, 7, 7]))
print(Solution().getDescentPeriods(prices=[8, 6, 7, 7]))
