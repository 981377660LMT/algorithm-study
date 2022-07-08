from collections import defaultdict
from fractions import Fraction
from typing import List

# 求最多有多少个点在同一条直线上。


class Solution:
    def maxPoints(self, points: List[List[int]]) -> int:
        """两点确定一条直线"""
        n, res = len(points), 1
        for i in range(n):
            x1, y1 = points[i]
            slopeGroup = defaultdict(list)  # 点i加斜率唯一确定一条直线
            for j in range(i + 1, n):
                x2, y2 = points[j]
                if x1 == x2:
                    slopeGroup[None].append(j)
                else:
                    slope = Fraction(y2 - y1, x2 - x1)  # 可以用元组+gcd代替
                    slopeGroup[slope].append(j)

            for slope, group in slopeGroup.items():
                res = max(res, len(group) + 1)

        return res


if __name__ == "__main__":
    print(Solution().maxPoints([[1, 1], [2, 2], [3, 3]]))
    print(Fraction(1, 0))  # ZeroDivisionError
