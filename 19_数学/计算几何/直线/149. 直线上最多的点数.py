from collections import defaultdict
from decimal import Decimal
from fractions import Fraction
from math import gcd
from typing import List, Tuple

# !求最多有多少个点在同一条直线上。
INF = int(1e18)


def calSlope1(x1: int, y1: int, x2: int, y2: int) -> Tuple[int, int]:
    """直线斜率"""
    if x2 == x1:
        return (INF, INF)
    gcd_ = gcd(x2 - x1, y2 - y1)
    a, b = (y2 - y1) // gcd_, (x2 - x1) // gcd_
    if a == 0:
        return (0, b if b > 0 else -b)
    elif a < 0:
        return (-a, -b)
    else:
        return (a, b)


class Solution:
    def maxPoints(self, points: List[List[int]]) -> int:
        """两点确定一条直线"""
        n, res = len(points), 1
        for i in range(n):
            x1, y1 = points[i]
            slopeCounter = defaultdict(int)  # 点i加斜率唯一确定一条直线
            for j in range(i + 1, n):
                x2, y2 = points[j]
                slope = calSlope1(x1, y1, x2, y2)
                print(slope)
                slopeCounter[slope] += 1

            for count in slopeCounter.values():
                res = max(res, count + 1)

        return res


if __name__ == "__main__":
    print(Solution().maxPoints([[1, 1], [2, 2], [3, 3]]))
    print(Solution().maxPoints([[2, 3], [3, 3], [-5, 3]]))
