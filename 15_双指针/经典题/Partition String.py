# 尽可能多的分割字符串，使得每种字符都只存在于一个子串中
from typing import List


class Solution:
    def solve(self, s: str) -> List[int]:
        last = {char: index for index, char in enumerate(s)}
        res = []
        lower, upper = 0, 0

        for index, char in enumerate(s):
            upper = max(upper, last[char])
            if index == upper:
                res.append(upper - lower + 1)
                lower = upper = index + 1

        return res


print(Solution().solve(s="cocoplaydae"))
