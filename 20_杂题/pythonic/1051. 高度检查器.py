from typing import List


class Solution:
    def heightChecker(self, heights: List[int]) -> int:
        return sum(n1 != n2 for n1, n2 in zip(sorted(heights), heights))


