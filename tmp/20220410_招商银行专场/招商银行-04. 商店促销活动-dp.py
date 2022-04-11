import fractions
from typing import List
from functools import lru_cache


MOD = int(1e9 + 7)
INF = int(1e20)

# A:购买商品满三件及以上商品可以打七折（向下取整）
# B:每购买三件商品，可以免去其中价格最低的一件商品的价格
# 假如需要编号 0 ～ N-1 的商品各买一件，请你求出最少的花费。
# n<=1e5


# 0.分情况dfs
# A如果触发了打折，那么是所有商品都7折，如果没有触发打折那么都是原价。
# A触发打折和A没触发打折要分两种情况讨论，所以有两个dfs

# 1.dfs前贪心
# 怎么让B的收益最大呢？ 可以一开始让B从大到小排序。这样每次免费的B一定是最贵的。
# 当B要变成3的时候，可以触发一次免费，然后把B变成0。

# 2.减少无用的状态数
# 注意买A商品应该用min(ca+1,3)而不是ca+1，来减少状态数
# 这是因为4，5，6...都没用，因为最后只需要判断买的A是否大于等于3

# 3. 使用Fraction来避免浮点数四则运算时产生的误差
# 有时候将浮点数或者Decimal作为Fraction实例的初始化数据可能会遇到舍入误差的问题
# costA *= fractions.Fraction(7, 10)  # 用分数类来计算,不能写*0.7/或者乘以7，最后除以10
print(0.1 + 0.2)
print(float((fractions.Fraction(1, 10) + fractions.Fraction(2, 10))))


DISCOUNT = fractions.Fraction(7, 10)
DISCOUNT = fractions.Fraction(0.7).limit_denominator()


class Solution:
    def goShopping(self, priceA: List[int], priceB: List[int]) -> int:
        """时间复杂度O(n*4*3)"""

        @lru_cache(None)
        def dfs1(index: int, ca: int, cb: int) -> float:
            """A购买三个或三个以上"""
            if index == n:
                return 0 if ca >= 3 else int(1e20)

            costA, costB = goods[index]
            costA *= DISCOUNT  # 用分数类来计算,不能写*0.7/或者打折的乘以7，不打折的乘以10，最后除以10
            if cb == 2:
                return min(dfs1(index + 1, ca, 0), dfs1(index + 1, min(3, ca + 1), cb) + costA)
            return min(
                dfs1(index + 1, ca, cb + 1) + costB, dfs1(index + 1, min(3, ca + 1), cb) + costA
            )

        @lru_cache(None)
        def dfs2(index: int, ca: int, cb: int) -> float:
            """A购买少于三个"""
            if ca > 2:
                return int(1e20)
            if index == n:
                return 0

            costA, costB = goods[index]
            if ca == 2:
                if cb == 2:
                    return dfs2(index + 1, ca, 0)
                else:
                    return dfs2(index + 1, ca, cb + 1) + costB
            if cb == 2:
                return min(dfs2(index + 1, min(3, ca + 1), cb) + costA, dfs2(index + 1, ca, 0))
            return min(
                dfs2(index + 1, ca, cb + 1) + costB, dfs2(index + 1, min(3, ca + 1), cb) + costA
            )

        n = len(priceA)
        goods = sorted(zip(priceA, priceB), key=lambda x: x[1], reverse=True)

        res1 = int(dfs1(0, 0, 0))
        dfs1.cache_clear()
        res2 = int(dfs2(0, 0, 0))
        dfs2.cache_clear()
        return min(res1, res2)


print(Solution().goShopping(priceA=[3, 13, 5, 12], priceB=[28, 12, 20, 7]))

