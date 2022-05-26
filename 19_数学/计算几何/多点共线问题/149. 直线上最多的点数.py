from collections import defaultdict
from fractions import Fraction
from typing import List

# 求最多有多少个点在同一条直线上。
class Solution:
    def maxPoints(self, points: List[List[int]]) -> int:
        n, res = len(points), 1
        for i in range(n):
            x1, y1 = points[i]
            counter = defaultdict(int)
            for j in range(i + 1, n):
                x2, y2 = points[j]
                if x1 == x2:
                    counter[(0, 0)] += 1
                else:
                    slope = Fraction(y2 - y1, x2 - x1)
                    counter[slope] += 1
            res = max(res, max(counter.values(), default=0) + 1)

        return res

