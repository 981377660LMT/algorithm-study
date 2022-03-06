import string
from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def cellsInRange(self, s: str) -> List[str]:
        words = s.split(':')
        nums = [int(n[1:]) for n in words]
        chars = [n[0] for n in words]
        min_ = min(nums)
        max_ = max(nums)
        res = []
        for char in string.ascii_uppercase:
            if chars[0] <= char <= chars[-1]:
                for i in range(min_, max_ + 1):
                    res.append(char + str(i))
        return res


["U7", "U8", "U9", "V7", "V8", "V9", "W7", "W8", "W9", "X7", "X8", "X9"]

