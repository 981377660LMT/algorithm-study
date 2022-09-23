from fractions import Fraction
from itertools import groupby, pairwise
from math import gcd
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# 计算斜率(x2!=x1)时):
# 1. Fraction((y2 - y1) / (x2 - x1)))
# 2. 元组 gcd_ = gcd(x2 - x1, y2 - y1)
# ((y2 - y1) // gcd_, (x2 - x1) // gcd_)


def calSlope(x1: int, y1: int, x2: int, y2: int) -> Tuple[int, int]:
    """直线斜率"""
    if x2 == x1:
        return (INF, INF)
    gcd_ = gcd(x2 - x1, y2 - y1)
    return ((y2 - y1) // gcd_, (x2 - x1) // gcd_)


class Solution:
    def minimumLines(self, stockPrices: List[List[int]]) -> int:
        # 注意原数据没有排序 需要sort一下
        slopes = (
            Fraction((y2 - y1), (x2 - x1)) for (x1, y1), (x2, y2) in pairwise(sorted(stockPrices))
        )
        return len(list(groupby(slopes)))

    def minimumLines2(self, stockPrices: List[List[int]]) -> int:
        stockPrices.sort()
        slopes = []
        for (x1, y1), (x2, y2) in pairwise(stockPrices):
            slopes.append(calSlope(x1, y1, x2, y2))
        groups = [[char, len(list(group))] for char, group in groupby(slopes)]
        return len(groups)


print(Solution().minimumLines([[1, 1], [500000000, 499999999], [1000000000, 999999998]]))
