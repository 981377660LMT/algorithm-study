from itertools import accumulate
from typing import List


class Solution:
    def isCovered(self, ranges: List[List[int]], left: int, right: int) -> bool:
        diff = [0] * 100
        for l, r in ranges:
            diff[l] += 1
            diff[r + 1] -= 1
        diff = list(accumulate(diff))
        return all(diff[i] > 0 for i in range(left, right + 1))

