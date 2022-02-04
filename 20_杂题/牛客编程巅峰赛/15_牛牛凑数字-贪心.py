# 从这些数字里面购买，并且凑到最大的数字带回家。
from typing import List


class Solution:
    def buyNumber(self, money: int, prices: List[int]) -> str:
        """从这些数字里面购买，并且凑到最大的数字带回家。"""
        """计算最便宜的数字，看最多买多少个，剩下的钱用于从高位到低位替换"""
        minPrice = min(prices)
        if minPrice > money:
            return "-1"

        count = money // minPrice
        remain = money - minPrice * count
        digit = next(i + 1 for i, v in reversed(list(enumerate(prices))) if v == minPrice)
        res = [str(digit)] * count

        index = 0
        for num in range(9, digit, -1):
            while index < len(res) and remain >= (prices[num - 1] - prices[digit - 1]):
                remain -= prices[num - 1] - prices[digit - 1]
                res[index] = str(num)
                index += 1
        return ''.join(res)


# print(Solution().buyNumber(5, [5, 4, 3, 2, 1, 2, 3, 4, 5]))
# print(Solution().buyNumber(5, [2, 3, 24, 32, 31, 14, 15, 7, 9]))
print(Solution().buyNumber(30274, [11010, 5050, 29464, 15646, 18587, 17315, 11139, 5293, 13249]))
