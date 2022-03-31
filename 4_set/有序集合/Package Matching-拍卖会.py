# 1e5
from typing import List

from sortedcontainers import SortedList


class Solution:
    def solve(self, sales: List[List[int]], buyers: List[List[int]]) -> int:
        """return the maximum number of packages that can be bought."""
        sales.sort()
        buyers.sort()

        res = index = 0
        offers = SortedList()

        # 对每一个物品，添加所有的销售人
        for day, price in sales:
            while index < len(buyers) and buyers[index][0] <= day:
                offers.add(buyers[index][1])
                index += 1

            select = offers.bisect_left(price)
            if select != len(offers):
                offers.pop(select)
                res += 1

        return res


print(
    Solution().solve(
        sales=[[0, 2], [0, 2], [0, 3], [1, 1], [1, 2], [3, 1]], buyers=[[0, 1], [0, 3], [1, 2]]
    )
)

# [day, price]
# day表示商品在day天上新
# [payday, amount]
# payday表示在payday以及之后都可以付款
