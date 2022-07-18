# 高度数组中的最大值对应的就是可重叠最长重复子串
from SA import useSA


class Solution:
    def longestDupSubstring(self, s: str) -> str:
        sa, _rank, height = useSA(list(map(ord, s)))
        maxI, max_ = 0, 0
        for i, h in enumerate(height):
            if h > max_:
                max_ = h
                maxI = i
        return s[sa[maxI] : sa[maxI] + max_]
