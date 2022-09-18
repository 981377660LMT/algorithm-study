"""6187. 完成所有交易的初始最少钱数"""

from typing import List


class Solution:
    def minimumMoney(self, transactions: List[List[int]]) -> int:
        """
        transactions[i] = [costi, cashbacki]
        为了完成交易 i ,money >= costi 这个条件必须为真。
        执行交易后，你的钱数 money 变成 money - costi + cashbacki 。
        请你返回 任意一种 交易顺序下，你都能完成所有交易的最少钱数 money 是多少

        贪心排序:
        1. 亏钱的排前面
        2. 亏钱的中, cashback 小的排前面 (不给发展机会)
        3. 赚钱的中, cost 大的排前面 (拦路虎)
        """
        bad, good = [], []
        for cost, cashback in transactions:
            actual = cost - cashback
            if actual > 0:
                bad.append((cost, cashback))
            else:
                good.append((cost, cashback))

        bad.sort(key=lambda x: x[1])
        good.sort(key=lambda x: x[0], reverse=True)
        nums = bad + good

        res, cur = 0, 0
        for cost, cashback in nums:
            if cur < cost:
                res += cost - cur
                cur = cost
            cur -= cost
            cur += cashback
        return res


print(Solution().minimumMoney(transactions=[[2, 1], [5, 0], [4, 2]]))
print(Solution().minimumMoney(transactions=[[3, 0], [0, 3]]))
# 10
# 3
# 24
print(
    Solution().minimumMoney(
        transactions=[
            [3, 9],
            [0, 4],
            [7, 10],
            [3, 5],
            [0, 9],
            [9, 3],
            [7, 4],
            [0, 0],
            [3, 3],
            [8, 0],
        ]
    )
)
# 24
