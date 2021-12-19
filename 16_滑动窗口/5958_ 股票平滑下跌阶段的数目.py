from typing import List

INF = 0x3FFFFFFF


class Solution:
    def getDescentPeriods(self, prices: List[int]) -> int:
        if len(prices) == 1:
            return 1
        res = 0
        left, right = 0, 0
        while right < len(prices):
            while right + 1 < len(prices) and prices[right] == prices[right + 1] + 1:
                right += 1
            diff = right - left + 1
            res += diff * (diff + 1) // 2
            right += 1
            left = right

        return res


print(Solution().getDescentPeriods(prices=[3, 2, 1, 4]))
print(Solution().getDescentPeriods(prices=[8, 6, 7, 7]))
print(Solution().getDescentPeriods(prices=[8, 6, 7, 7]))
