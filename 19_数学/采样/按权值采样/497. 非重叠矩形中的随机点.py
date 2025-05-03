from bisect import bisect_right
from random import randrange
from typing import List


class Solution:
    def __init__(self, rects: List[List[int]]):
        self.rects = rects
        self.presum = [0]
        for a, b, x, y in rects:
            self.presum.append(self.presum[-1] + (x - a + 1) * (y - b + 1))

    def pick(self) -> List[int]:
        k = randrange(self.presum[-1])
        rectIndex = bisect_right(self.presum, k) - 1
        a, b, _, y = self.rects[rectIndex]
        da, db = divmod(k - self.presum[rectIndex], y - b + 1)
        return [a + da, b + db]
