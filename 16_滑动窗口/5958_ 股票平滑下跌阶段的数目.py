from typing import List


class Solution:
    def getDescentPeriods(self, prices: List[int]) -> int:
        n = len(prices)
        res = 1
        count = 1
        for i in range(1, n):
            if prices[i] == prices[i - 1] - 1:
                count += 1
            else:
                count = 1
            res += count
        return res


print(Solution().getDescentPeriods(prices=[3, 2, 1, 4]))
print(Solution().getDescentPeriods(prices=[8, 6, 7, 7]))
print(Solution().getDescentPeriods(prices=[8, 6, 7, 7]))
