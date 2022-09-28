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
        for count in ratios.values():
            res += count * (count - 1) // 2
        return res

    def interchangeableRectangles2(self, rectangles: List[List[int]]) -> int:
        """一遍遍历"""
        res, ratio = 0, dict()
        for w, h in rectangles:
            gcd_ = gcd(w, h)
            w, h = w // gcd_, h // gcd_
            res += ratio.get((w, h), 0)
            ratio[(w, h)] = ratio.get((w, h), 0) + 1
        return res


print(Solution().interchangeableRectangles(rectangles=[[4, 8], [3, 6], [10, 20], [15, 30]]))
