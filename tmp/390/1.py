from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)


# 给你一个字符串 s ，请找出满足每个字符最多出现两次的最长子字符串，并返回该子字符串的 最大 长度。
class Solution:
    def maximumLengthSubstring(self, s: str) -> int:
        res = 0
        for i in range(len(s)):
            for j in range(i + 1, len(s)):
                counter = Counter(s[i : j + 1])
                if all([v <= 2 for v in counter.values()]):
                    res = max(res, j - i + 1)
        return res
