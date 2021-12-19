# 4x+2y=a x+y=b
from typing import List


class Solution:
    def numOfBurgers(self, tomatoSlices: int, cheeseSlices: int) -> List[int]:
        s = (4 * cheeseSlices - tomatoSlices) / 2
        b = (tomatoSlices - 2 * cheeseSlices) / 2
        if s >= 0 and b >= 0 and s == int(s) and b == int(b):
            return [int(b), int(s)]
        return []
