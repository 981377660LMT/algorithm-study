"""
!采购水果的最大美味度
总费用不能超过maxAmount 可以使用maxCoupons次优惠券来将水果的价格减半(向下取整)
求最大美味度


# n == price.length == tastiness.length
# 1 <= n <= 100
# 0 <= price[i], tastiness[i], maxAmount <= 1000
# 0 <= maxCoupons <= 5
多维费用的01背包问题
"""

from functools import lru_cache
from typing import List

INF = int(1e18)


class Solution:
    def maxTastiness(
        self, price: List[int], tastiness: List[int], maxAmount: int, maxCoupons: int
    ) -> int:
        @lru_cache(None)
        def dfs(index: int, money: int, coupons: int) -> int:
            if money < 0 or coupons < 0:
                return -INF
            if index == n or money == 0:
                return 0

            # jump
            res = dfs(index + 1, money, coupons)

            cost, score = price[index], tastiness[index]
            # buy without coupon
            if cost <= money:
                cand = dfs(index + 1, money - cost, coupons) + score
                res = cand if cand > res else res

            # buy with coupon
            if cost // 2 <= money:
                cand = dfs(index + 1, money - cost // 2, coupons - 1) + score
                res = cand if cand > res else res

            return res

        n = len(price)
        res = dfs(0, maxAmount, maxCoupons)
        dfs.cache_clear()
        return res
