from itertools import groupby


class Solution:
    def makeFancyString(self, s: str) -> str:
        # 至多取两个
        return ''.join(next(g) + next(g, '') for _, g in groupby(s))
