from typing import List
from collections import Counter
from math import gcd


class Solution:
    def interchangeableRectangles(self, rectangles: List[List[int]]) -> int:
        ratios = Counter()
        for x, y in rectangles:
            g = gcd(x, y)
            ratios[(x / g, y / g)] += 1

        res = 0
        for val in ratios.values():
            res += val * (val - 1) // 2
        return res


print(Solution().interchangeableRectangles(rectangles=[[4, 8], [3, 6], [10, 20], [15, 30]]))
