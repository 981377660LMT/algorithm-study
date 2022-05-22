from fractions import Fraction
from itertools import groupby, pairwise
from math import gcd
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# 计算斜率:
# 1. Fraction((y2 - y1) / (x2 - x1)))
# 2. 元组 gcd_ = gcd(x2 - x1, y2 - y1)
# ((y2 - y1) // gcd_, (x2 - x1) // gcd_)


class Solution:
    def minimumLines(self, stockPrices: List[List[int]]) -> int:
        # 注意原数据没有排序
        slopes = (
            Fraction((y2 - y1), (x2 - x1)) for (x1, y1), (x2, y2) in pairwise(sorted(stockPrices))
        )
        return len(list(groupby(slopes)))

    def minimumLines2(self, stockPrices: List[List[int]]) -> int:
        stockPrices.sort()
        slopes = []
        for (x1, y1), (x2, y2) in pairwise(stockPrices):
            gcd_ = gcd(x2 - x1, y2 - y1)
            slopes.append(((x2 - x1) // gcd_, (y2 - y1) // gcd_))
        groups = [[char, len(list(group))] for char, group in groupby(slopes)]
        return len(groups)


print(Solution().minimumLines([[1, 1], [500000000, 499999999], [1000000000, 999999998]]))
# 29
print(
    Solution().minimumLines(
        [
            [72, 98],
            [62, 27],
            [32, 7],
            [71, 4],
            [25, 19],
            [91, 30],
            [52, 73],
            [10, 9],
            [99, 71],
            [47, 22],
            [19, 30],
            [80, 63],
            [18, 15],
            [48, 17],
            [77, 16],
            [46, 27],
            [66, 87],
            [55, 84],
            [65, 38],
            [30, 9],
            [50, 42],
            [100, 60],
            [75, 73],
            [98, 53],
            [22, 80],
            [41, 61],
            [37, 47],
            [95, 8],
            [51, 81],
            [78, 79],
            [57, 95],
        ]
    )
)
