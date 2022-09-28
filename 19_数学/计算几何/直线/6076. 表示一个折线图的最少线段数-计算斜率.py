from decimal import Decimal
from itertools import groupby, pairwise
from math import gcd
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# 计算斜率:
# 1. Decimal
# 2. 元组


def calSlope1(x1: int, y1: int, x2: int, y2: int) -> Tuple[int, int]:
    """直线斜率"""
    if x2 == x1:
        return (INF, INF)
    gcd_ = gcd(x2 - x1, y2 - y1)  # !注意math.gcd总返回非负数
    a, b = (y2 - y1) // gcd_, (x2 - x1) // gcd_
    if a == 0:
        return (0, b if b > 0 else -b)
    elif a < 0:
        return (-a, -b)
    else:
        return (a, b)


def calSlope2(x1: int, y1: int, x2: int, y2: int):
    """直线斜率"""
    if x2 == x1:
        return INF
    return Decimal(y2 - y1) / Decimal(x2 - x1)


class Solution:
    def minimumLines(self, stockPrices: List[List[int]]) -> int:
        # 注意原数据没有排序 需要sort一下
        slopes = (
            Decimal(y2 - y1) / Decimal(x2 - x1)
            for (x1, y1), (x2, y2) in pairwise(sorted(stockPrices))
        )
        return len(list(groupby(slopes)))

    def minimumLines2(self, stockPrices: List[List[int]]) -> int:
        stockPrices.sort()
        slopes = []
        for (x1, y1), (x2, y2) in pairwise(stockPrices):
            slopes.append(calSlope1(x1, y1, x2, y2))
        groups = [[char, len(list(group))] for char, group in groupby(slopes)]
        return len(groups)


print(Solution().minimumLines([[1, 1], [500000000, 499999999], [1000000000, 999999998]]))
