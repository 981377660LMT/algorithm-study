from typing import List


class Solution:
    def maxWidthOfVerticalArea(self, points: List[List[int]]) -> int:
        p = sorted(set([x for x, _ in points]))
        return max([b - a for a, b in zip(p, p[1:])], default=0)

